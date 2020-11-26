package gps_test

import (
	"testing"

	"github.com/andycondon/pathfinder/pkg/gps"
	"github.com/stretchr/testify/assert"
)

func TestFromGPRMC(t *testing.T) {
	tests := []struct {
		name     string
		sentence string
		reading  gps.Reading
	}{
		{
			name:     "no fix",
			sentence: `$GPRMC,215952.087,V,,,,,0.00,0.00,070180,,,N*44`,
			reading:  gps.Reading{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reading := gps.FromGPRMC(tc.sentence)

			assert.Equal(t, tc.reading, reading)
		})
	}
}
