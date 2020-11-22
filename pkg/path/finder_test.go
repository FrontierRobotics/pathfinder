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
