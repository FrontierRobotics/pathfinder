package path_test

import (
	"testing"

	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/motor"
	"github.com/andycondon/pathfinder/pkg/path"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name      string
		irReading ir.Reading
		driveCmd  motor.Command
	}{
		{
			name:      "move forward if all clear",
			irReading: ir.Reading{L: ir.ProximityClear, F: ir.ProximityClear, R: ir.ProximityClear},
			driveCmd:  motor.Command{M: motor.Forward, S: motor.Slow},
		},
		{
			name:      "rotate right if possible and obstacle to left",
			irReading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityFar, R: ir.ProximityFar},
			driveCmd:  motor.Command{M: motor.RotateRight, S: motor.Medium},
		},
		{
			name:      "rotate left if possible and obstacle to right",
			irReading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityNear},
			driveCmd:  motor.Command{M: motor.RotateLeft, S: motor.Medium},
		},
		{
			name:      "rotate left if front and right blocked",
			irReading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityNear, R: ir.ProximityNear},
			driveCmd:  motor.Command{M: motor.RotateLeft, S: motor.Medium},
		},
		{
			name:      "rotate right if everything blocked",
			irReading: ir.Reading{L: ir.ProximityNear, F: ir.ProximityNear, R: ir.ProximityNear},
			driveCmd:  motor.Command{M: motor.RotateRight, S: motor.Medium},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			irCh := make(chan ir.Reading)
			driverCh := make(chan motor.Command)

			f := path.Finder{IR: irCh, Drive: driverCh}

			go f.Find()
			irCh <- tc.irReading
			assert.Equal(t, tc.driveCmd, <-driverCh)
		})
	}
}
