package nav

import (
	"fmt"
	"math"
)

type Angle float64

const (
	Degrees Angle = 1
	Minutes       = Degrees / 60
	Seconds       = Degrees / 3600
	Radians       = math.Pi / 180
)

type Direction int

// A Speed represents the change in position
// as a float64 meters per second count.
type Speed float64

const (
	MetersPerSecond   Speed = 1
	KilometersPerHour       = 3.6
	Knots                   = 0.51444444444
)

type Position struct {
	Latitude  Angle
	Longitude Angle
}

func (p Position) String() string {
	return fmt.Sprintf("%f, %f", p.Latitude, p.Longitude)
}
