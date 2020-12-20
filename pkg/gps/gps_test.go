package gps_test

import (
	"testing"
	"time"

	"github.com/andycondon/pathfinder/pkg/gps"
	"github.com/andycondon/pathfinder/pkg/nav"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

// Times from the 80s might show up without a GPS fix. We'll just be dumb about the century.
var jan072080, _ = time.Parse(time.RFC3339, "2080-01-07T21:59:52.087Z")
var nov272020, _ = time.Parse(time.RFC3339, "2020-11-27T19:17:36.000Z")

// Useful for conversion checks: https://www.pgc.umn.edu/apps/convert/
func TestFromGPRMC(t *testing.T) {
	tests := []struct {
		name     string
		sentence string
		reading  gps.Reading
		hasErr   bool
	}{
		{
			name:     "no fix",
			sentence: `$GPRMC,215952.087,V,,,,,0.00,0.00,070180,,,N*44`,
			reading:  gps.Reading{Time: jan072080},
		},
		{
			name:     "fix",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,10448.5048,W,2.01,34.47,271120,,,A*48`,
			reading: gps.Reading{
				Time:     nov272020,
				Fix:      true,
				Speed:    1.0340333333244 * nav.MetersPerSecond,
				Position: s2.LatLngFromDegrees(41.18582333333333, -104.80841333333333),
			},
		},
		{
			name:     "fix down under",
			sentence: `$GPRMC,191736.000,A,2747.2403,S,13629.5662,E,2.01,34.47,271120,,,A*48`,
			reading: gps.Reading{
				Time:     nov272020,
				Fix:      true,
				Speed:    1.0340333333244 * nav.MetersPerSecond,
				Position: s2.LatLngFromDegrees(-27.787338333333334, 136.492770),
			},
		},
		{
			name:     "missing component",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,10448.5048,W,2.01,34.47`,
			hasErr:   true,
		},
		{
			name:     "empty time",
			sentence: `$GPRMC,,A,4111.1494,N,10448.5048,W,2.01,34.47,271120,,,A*48`,
			reading: gps.Reading{
				Fix:      true,
				Speed:    1.0340333333244 * nav.MetersPerSecond,
				Position: s2.LatLngFromDegrees(41.18582333333333, -104.80841333333333),
			},
		},
		{
			name:     "invalid time",
			sentence: `$GPRMC,i,A,4111.1494,N,10448.5048,W,2.01,34.47,271120,,,A*48`,
			hasErr:   true,
		},
		{
			name:     "empty date",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,10448.5048,W,2.01,34.47,,,,A*48`,
			reading: gps.Reading{
				Fix:      true,
				Speed:    1.0340333333244 * nav.MetersPerSecond,
				Position: s2.LatLngFromDegrees(41.18582333333333, -104.80841333333333),
			},
		},
		{
			name:     "invalid date",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,10448.5048,W,2.01,34.47,i,,,A*48`,
			hasErr:   true,
		},
		{
			name:     "invalid date month",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,10448.5048,W,2.01,34.47,27ii20,,,A*48`,
			hasErr:   true,
		},
		{
			name:     "empty speed",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,10448.5048,W,,34.47,271120,,,A*48`,
			reading: gps.Reading{
				Time:     nov272020,
				Fix:      true,
				Position: s2.LatLngFromDegrees(41.18582333333333, -104.80841333333333),
			},
		},
		{
			name:     "invalid speed",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,10448.5048,W,i,34.47,271120,,,A*48`,
			hasErr:   true,
		},
		{
			name:     "invalid latitue",
			sentence: `$GPRMC,191736.000,A,iii,N,10448.5048,W,2.01,34.47,271120,,,A*48`,
			hasErr:   true,
		},
		{
			name:     "invalid latitue minutes",
			sentence: `$GPRMC,191736.000,A,23i,N,10448.5048,W,2.01,34.47,271120,,,A*48`,
			hasErr:   true,
		},
		{
			name:     "invalid latitue length",
			sentence: `$GPRMC,191736.000,A,i,N,10448.5048,W,2.01,34.47,271120,,,A*48`,
			hasErr:   true,
		},
		{
			name:     "invalid longitude",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,iii,W,2.01,34.47,271120,,,A*48`,
			hasErr:   true,
		},
		{
			name:     "invalid longitude length",
			sentence: `$GPRMC,191736.000,A,4111.1494,N,i,W,2.01,34.47,271120,,,A*48`,
			hasErr:   true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reading, err := gps.FromGPRMC(tc.sentence)

			if tc.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.reading, reading)
			}
		})
	}
}

func TestReading_String(t *testing.T) {
	t.Run("no fix", func(t *testing.T) {
		r := gps.Reading{Time: jan072080}
		assert.Equal(t, `fix: none, time: 2080-01-07 21:59:52.087 +0000 UTC`, r.String())
	})

	t.Run("fix", func(t *testing.T) {
		r := gps.Reading{
			Time:     nov272020,
			Fix:      true,
			Speed:    1.0340333333244 * nav.MetersPerSecond,
			Position: s2.LatLngFromDegrees(41.1858233, -104.8084133),
		}
		assert.Equal(t, `fix: active, time: 2020-11-27 19:17:36 +0000 UTC, speed: 1.034033 m/s, position: [41.1858233, -104.8084133]`, r.String())
	})
}
