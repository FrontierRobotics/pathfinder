package arduino_test

import (
	"errors"
	"testing"

	"github.com/andycondon/pathfinder/pkg/arduino"
	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/stretchr/testify/assert"
)

func TestArduino_GetStatus(t *testing.T) {
	tests := []struct {
		name    string
		bytes   []byte
		reading ir.Reading
		hasErr  bool
		txErr   error
	}{
		{
			name:    "happy path",
			bytes:   []byte{0x20, 0x20, 0x20},
			reading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityFar},
		},
		{
			name:   "tx error",
			txErr:  errors.New("tx error"),
			hasErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			irArray := &ir.SensorArray{
				Left:    ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0x50},
				Forward: ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0x50},
				Right:   ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0x50},
			}
			tx := func(w []byte, r []byte) error {
				for i := range tc.bytes {
					r[i] = tc.bytes[i]
				}
				return tc.txErr
			}
			a := arduino.New(tx, irArray)

			status, err := a.GetStatus()

			if tc.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.reading, status.IR)
		})
	}
}
