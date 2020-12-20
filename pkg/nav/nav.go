package nav

import "github.com/golang/geo/s2"

// A Speed represents the change in position
// as a float64 meters per second count.
type Speed float64

type Distance float64

const (
	Meters     Distance = 1
	Kilometers          = Meters * 1000
)

func (d Distance) Meters() float64 { return float64(d) }

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
