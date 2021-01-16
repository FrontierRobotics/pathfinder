package nav_test

import (
	"testing"

	"github.com/andycondon/pathfinder/pkg/nav"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestDistance(t *testing.T) {
	var tests = []struct {
		name    string
		str     string
		in, out nav.Distance
	}{
		{
			name: "meters to kilometers",
			in:   36 * nav.Meters,
			out:  0.036 * nav.Kilometers,
			str:  "36.00 m",
		},
		{
			name: "kilometers to meters",
			in:   36 * nav.Kilometers,
			out:  36000 * nav.Meters,
			str:  "36.00 km",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.InDelta(t, tc.out.Meters(), tc.in.Meters(), 0.0001)
			assert.Equal(t, tc.str, tc.out.String())
		})
	}
}

func TestSpeed(t *testing.T) {
	var tests = []struct {
		name    string
		in, out nav.Speed
	}{
		{
			name: "mps to kph",
			in:   36 * nav.MetersPerSecond,
			out:  129.6 * nav.KilometersPerHour,
		},
		{
			name: "mph to kph",
			in:   75 * nav.MilesPerHour,
			out:  120.701 * nav.KilometersPerHour,
		},
		{
			name: "mph to mps",
			in:   75 * nav.MilesPerHour,
			out:  33.528 * nav.MetersPerSecond,
		},
		{
			name: "knots to mps",
			in:   75 * nav.Knots,
			out:  38.5833 * nav.MetersPerSecond,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.InDelta(t, tc.out.MetersPerSecond(), tc.in.MetersPerSecond(), 0.0001)
		})
	}
}

func TestRelativeBearing(t *testing.T) {
	tests := []struct {
		name    string
		h, b, r float64
	}{
		{name: "both in first quadrant, smaller heading", h: 10, b: 15, r: -5},
		{name: "both in first quadrant, larger heading", h: 20, b: 15, r: 5},
		{name: "both in fourth quadrant, smaller heading", h: 310, b: 315, r: -5},
		{name: "both in fourth quadrant, larger heading", h: 320, b: 315, r: 5},
		{name: "bearing in fourth quadrant, heading in first quadrant", h: 10, b: 350, r: 20},
		{name: "heading in fourth quadrant, bearing in first quadrant", h: 350, b: 10, r: -20},
		{name: "bearing in second quadrant, heading in third quadrant", h: 200, b: 150, r: 50},
		{name: "heading in second quadrant, bearing in third quadrant", h: 150, b: 200, r: -50},
		{name: "close enough", h: 10.01, b: 10, r: 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := s1.Angle(tc.h) * s1.Degree
			b := s1.Angle(tc.b) * s1.Degree

			expected := s1.Angle(tc.r) * s1.Degree
			actual := nav.RelativeBearing(h, b, 0.1*s1.Degree)

			assert.InDelta(t, expected.Degrees(), actual.Degrees(), 0.0001)
		})
	}
}

// Lat/Lon values from Google Maps
// Distance expected values from Google maps
func TestDistanceBetweenPositions(t *testing.T) {
	var tests = []struct {
		name   string
		p1, p2 s2.LatLng
		dist   nav.Distance
		delta  nav.Distance
	}{
		{
			name:  "across the house",
			p1:    s2.LatLngFromDegrees(41.18567319147294, -104.80829286889617),
			p2:    s2.LatLngFromDegrees(41.18567113019101, -104.80852907096639),
			dist:  19.49 * nav.Meters,
			delta: 0.3 * nav.Meters,
		},
		{
			name:  "to the middle of a field",
			p1:    s2.LatLngFromDegrees(41.18579691821355, -104.80842287569891),
			p2:    s2.LatLngFromDegrees(41.169644780072325, -104.80027094457998),
			dist:  1.921 * nav.Kilometers,
			delta: 0.3 * nav.Meters,
		},
		{
			name:  "across town",
			p1:    s2.LatLngFromDegrees(41.15548155514516, -104.6532153220288),
			p2:    s2.LatLngFromDegrees(41.176720961200395, -104.8469773033338),
			dist:  16.31 * nav.Kilometers,
			delta: 100 * nav.Meters,
		},
		{
			name:  "across state",
			p1:    s2.LatLngFromDegrees(41.020940972893584, -104.05014439015459),
			p2:    s2.LatLngFromDegrees(41.01566725945257, -111.05755219078006),
			dist:  586.54 * nav.Kilometers,
			delta: 2 * nav.Kilometers,
		},
		{
			name:  "to Hawaii",
			p1:    s2.LatLngFromDegrees(41.18905770370119, -104.80451497634355),
			p2:    s2.LatLngFromDegrees(19.68913610610713, -155.45520158844647),
			dist:  5326.49 * nav.Kilometers,
			delta: 10 * nav.Kilometers,
		},
		{
			name:  "to Antarctica",
			p1:    s2.LatLngFromDegrees(41.18905770370119, -104.80451497634355),
			p2:    s2.LatLngFromDegrees(-69.99939496893386, 86.52316700898501),
			dist:  16720.95 * nav.Kilometers,
			delta: 30 * nav.Kilometers,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := nav.DistanceBetweenPositions(tc.p1, tc.p2)
			assert.InDelta(t, d.Meters(), tc.dist.Meters(), tc.delta.Meters())
		})
	}
}

func TestInitialBearing(t *testing.T) {
	var tests = []struct {
		name   string
		p1, p2 s2.LatLng
		b      s1.Angle
	}{
		{
			name: "across the house",
			p1:   s2.LatLngFromDegrees(41.18567319147294, -104.80829286889617),
			p2:   s2.LatLngFromDegrees(41.18567113019101, -104.80852907096639),
			b:    s1.Angle(4.7008),
		},
		{
			name: "down the driveway",
			p1:   s2.LatLngFromDegrees(41.185579119692925, -104.80847684264579),
			p2:   s2.LatLngFromDegrees(41.185493329876394, -104.80846343169397),
			b:    s1.Angle(3.0245),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := nav.InitialBearing(tc.p1, tc.p2)
			assert.InDelta(t, b.Radians(), tc.b.Radians(), .01)
		})
	}
}
