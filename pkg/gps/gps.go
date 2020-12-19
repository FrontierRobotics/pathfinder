package gps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andycondon/pathfinder/pkg/nav"
)

type Reading struct {
	Fix      bool
	Time     time.Time
	Speed    nav.Speed
	Position nav.Position
}

func (r Reading) String() string {
	if r.Fix {
		return fmt.Sprintf("fix: active, time: %s, speed: %f m/s, position: %s", r.Time.String(), r.Speed, r.Position.String())
	}
	return fmt.Sprintf("fix: none, time: %s", r.Time.String())
}

// GPRMC Sentence
// see: https://learn.adafruit.com/adafruit-ultimate-gps/direct-computer-wiring
// $GPRMC,191736.000,A,4111.1494,N,10448.5048,W,0.01,34.47,271120,,,A*48
//      0,         1,2,        3,4,         5,6,   7,    8,     9,,,12
// 1: Timestamp in UTC - HHMMSS.mmm
// 2: Status, V = void, A = active (locked)
// 3: Latitude angle
// 4: Latitude direction
// 5: Longitude angle
// 6: Longitude direction
// 7: Ground speed in knots
// 8: Tracking angle - approximate direction we're headed in
// 9: Datestamp - DDMMYY
// 12: Checksum

// We could derive this from time.Now, but if we're setting our clock from GPS that defeats the point
const currentCentury = 2000

// See http://aprs.gids.nl/nmea/#rmc
func FromGPRMC(sentence string) (Reading, error) {
	parts := strings.Split(sentence, ",")
	if len(parts) != 13 {
		return Reading{}, errors.New("sentence must have 13 parts")
	}

	t, err := timeFromSentence(parts[1], parts[9])
	if err != nil {
		return Reading{}, err
	}

	speed, err := speedFromSentence(parts[7])
	if err != nil {
		return Reading{}, err
	}

	position, err := postitionFromSentence(parts[3], parts[4], parts[5], parts[6])
	if err != nil {
		return Reading{}, err
	}

	return Reading{
		Time:     t,
		Fix:      parts[2] == "A",
		Speed:    speed,
		Position: position,
	}, nil
}

func timeFromSentence(t, d string) (ts time.Time, err error) {
	if t == "" || d == "" {
		return time.Time{}, nil
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("error parsing time: %v", r))
		}
	}()
	var (
		dayOfMonth = toInt(d[0:2])
		month      = time.Month(toInt(d[2:4]))
		year       = toInt(d[4:]) + currentCentury

		hour        = toInt(t[0:2])
		minute      = toInt(t[2:4])
		seconds     = toInt(t[4:6])
		nanoseconds = toInt(t[7:]) * 1000000
	)

	return time.Date(year, month, dayOfMonth, hour, minute, seconds, nanoseconds, time.UTC), err
}

func speedFromSentence(s string) (nav.Speed, error) {
	if s == "" {
		return 0, nil
	}
	t, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return nav.Speed(nav.Knots * t), nil
}

func postitionFromSentence(lat, latDir, lon, lonDir string) (nav.Position, error) {
	if lat == "" || latDir == "" || lon == "" || lonDir == "" {
		return nav.Position{}, nil
	}
	if len(lat) < 2 {
		return nav.Position{}, errors.New("invalid format for latitude")
	}
	if len(lon) < 3 {
		return nav.Position{}, errors.New("invalid format for longitude")
	}
	latitude, err := angleFromSentence(lat, 2)
	if err != nil {
		return nav.Position{}, err
	}
	longitude, err := angleFromSentence(lon, 3)
	if err != nil {
		return nav.Position{}, err
	}
	if latDir == "S" {
		latitude = latitude * -1
	}
	if lonDir == "W" {
		longitude = longitude * -1
	}
	return nav.Position{Latitude: latitude, Longitude: longitude}, nil
}

func angleFromSentence(a string, minLen int) (nav.Angle, error) {
	deg, err := strconv.ParseFloat(a[0:minLen], 64)
	if err != nil {
		return 0, err
	}
	min, err := strconv.ParseFloat(a[minLen:], 64)
	if err != nil {
		return 0, err
	}
	degrees := nav.Angle(deg) * nav.Degrees
	minutes := nav.Angle(min) * nav.Minutes
	return degrees + minutes, nil
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
