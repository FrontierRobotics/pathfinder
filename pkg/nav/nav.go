package nav

import (
	"fmt"
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// A Speed represents the change in position
// as a float64 meters per second count.
type Speed float64

type Distance float64

const (
	Meters     Distance = 1
	Kilometers          = Meters * 1000
)

func (d Distance) Meters() float64 { return float64(d) }

func (d Distance) Kilometers() float64 { return float64(d / Kilometers) }

func (d Distance) String() string {
	if d < Kilometers {
		return fmt.Sprintf("%0.2f m", d)
	}
	return fmt.Sprintf("%0.2f km", d.Kilometers())
}

const (
	MetersPerSecond   Speed = 1
	KilometersPerHour       = MetersPerSecond / 3.6
	MilesPerHour            = 0.44704
	Knots                   = 0.51444444444
)

func (s Speed) MetersPerSecond() float64 { return float64(s) }

const R = 6371 * Kilometers

func DistanceBetweenPositions(p1, p2 s2.LatLng) Distance {
	c := p1.Distance(p2)
	return Distance(R.Meters() * c.Radians())
}

// RelativeBearing returns the shortest difference between a given
// heading h, and bearing b.
// It will return a zero angle if the difference does not exceed
// a tolerance t.
// See https://stackoverflow.com/a/7869457
func RelativeBearing(h, b, t s1.Angle) s1.Angle {
	r := h - b
	if r > math.Pi {
		r = r - (2 * math.Pi)
	} else if r < (-1 * math.Pi) {
		r = r + (2 * math.Pi)
	}
	if r.Abs() < t {
		return 0
	}
	return r
}

/* Proud of this, but not as good as above
r := h - b
if r.Abs() > math.Pi {
	if r < 0 {
		return (h + 2*math.Pi) - b
	} else {
		return ((b + 2*math.Pi) - h) * -1
	}
}
*/

// sourced from https://github.com/google/s2geometry/blob/master/src/s2/s2earth.cc
// sourced from https://www.movable-type.co.uk/scripts
// Formula: θ = atan2( sin Δλ ⋅ cos φ2 , cos φ1 ⋅ sin φ2 − sin φ1 ⋅ cos φ2 ⋅ cos Δλ )
// where    φ1,λ1 is the start point, φ2,λ2 the end point (Δλ is the difference in longitude)
func InitialBearing(a, b s2.LatLng) s1.Angle {
	var (
		lat1    = a.Lat.Radians()
		cosLat2 = math.Cos(b.Lat.Radians())
		latDiff = b.Lat.Radians() - a.Lat.Radians()
		lngDiff = b.Lng.Radians() - a.Lng.Radians()

		x     = math.Sin(latDiff) + math.Sin(lat1)*cosLat2*2*haversine(lngDiff)
		y     = math.Sin(lngDiff) * cosLat2
		theta = math.Atan2(y, x)
	)
	if theta < 0 {
		return s1.Angle(theta + 2*math.Pi)
	}
	return s1.Angle(theta)
}

// http://en.wikipedia.org/wiki/Haversine_formula
// Haversine(x) has very good numerical stability around zero.
// Haversine(x) == (1-cos(x))/2 == sin(x/2)^2; must be implemented with the
// second form to reap the numerical benefits.
func haversine(radians float64) float64 {
	sinHalf := math.Sin(radians / 2)
	return sinHalf * sinHalf
}
