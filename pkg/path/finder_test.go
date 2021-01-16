package path

import (
	"testing"

	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/motor"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestFinder_Avoidance(t *testing.T) {
	tests := []struct {
		name      string
		irReading ir.Reading
		driveCmd  motor.Command
	}{
		{
			name:      "move fast if all clear",
			irReading: ir.Reading{L: ir.ProximityClear, F: ir.ProximityClear, R: ir.ProximityClear},
			driveCmd:  motor.Command{M: motor.Forward, S: motor.Fast},
		},
		{
			name:      "move forward if all far",
			irReading: ir.Reading{L: ir.ProximityFar, F: ir.ProximityFar, R: ir.ProximityFar},
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
			_, cmd := avoid(tc.irReading)
			assert.Equal(t, tc.driveCmd, cmd)
		})
	}
}

func TestFind(t *testing.T) {
	before := func(waypoints ...s2.LatLng) (chan bool, chan s2.LatLng, chan s1.Angle, chan motor.Command, Finder) {
		var (
			fix      = make(chan bool)
			position = make(chan s2.LatLng)
			heading  = make(chan s1.Angle)
			drive    = make(chan motor.Command, 3)
			f        = Finder{
				GPSfix:   fix,
				Position: position, Heading: heading, Drive: drive,
				Waypoints:        waypoints,
				BearingTolerance: 1 * s1.Degree,
			}
		)
		return fix, position, heading, drive, f
	}
	t.Run("changes heading to match bearing and heads to waypoint", func(t *testing.T) {
		fix, position, heading, drive, f := before(
			s2.LatLngFromDegrees(41.18567113019101, -104.80852907096639),
		)
		go f.Find()

		//assert.Equal(t, motor.Command{M: motor.Park, S: motor.Stop}, <-drive) // waiting for fix
		fix <- true
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		position <- s2.LatLngFromDegrees(41.18567319147294, -104.80829286889617)
		assert.Equal(t, motor.Command{M: motor.RotateLeft, S: motor.Medium}, <-drive)

		heading <- 269 * s1.Degree
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)
	})
	t.Run("stops at waypoint", func(t *testing.T) {
		fix, position, heading, drive, f := before(
			s2.LatLngFromDegrees(41.18567113019101, -104.80852907096639),
		)
		go f.Find()

		fix <- true
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		position <- s2.LatLngFromDegrees(41.185579119692925, -104.80847684264579)
		assert.Equal(t, motor.Command{M: motor.RotateLeft, S: motor.Medium}, <-drive)

		heading <- 336 * s1.Degree
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		position <- s2.LatLngFromDegrees(41.18567113019101, -104.80852907096639)
		assert.Equal(t, motor.Command{M: motor.Park, S: motor.Stop}, <-drive)
	})
	t.Run("stops if loses fix", func(t *testing.T) {
		fix, position, heading, drive, f := before(
			s2.LatLngFromDegrees(41.18567113019101, -104.80852907096639),
		)
		go f.Find()

		fix <- true
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		position <- s2.LatLngFromDegrees(41.185579119692925, -104.80847684264579)
		assert.Equal(t, motor.Command{M: motor.RotateLeft, S: motor.Medium}, <-drive)

		heading <- 336 * s1.Degree
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		fix <- false
		assert.Equal(t, motor.Command{M: motor.Park, S: motor.Stop}, <-drive)
	})
	t.Run("follows a mission", func(t *testing.T) {
		/*
			// start
			41.185579119692925, -104.80847684264579
			// midpoint 1
			41.18553067368046, -104.80847416049401
			// waypoint 1
			41.185493329876394, -104.80846343169397
			// midpoint 2
			41.18550645031859, -104.80825556017193
			// waypoint 2
			41.18551556476105, -104.80804562107231
		*/
		fix, position, heading, drive, f := before(
			// waypoint 1
			s2.LatLngFromDegrees(41.185493329876394, -104.80846343169397),
			// waypoint 2
			s2.LatLngFromDegrees(41.18551556476105, -104.80804562107231),
		)
		go f.Find()

		fix <- true
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		// start
		position <- s2.LatLngFromDegrees(41.185579119692925, -104.80847684264579)
		assert.Equal(t, motor.Command{M: motor.RotateRight, S: motor.Medium}, <-drive) // correcting heading

		heading <- s1.Degree * 229                                                    // first heading reading, but in wrong direction
		assert.Equal(t, motor.Command{M: motor.RotateLeft, S: motor.Medium}, <-drive) // correcting heading

		heading <- s1.Degree * 173 // heading now corrected to match bearing
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		// midpoint 1
		position <- s2.LatLngFromDegrees(41.18553067368046, -104.80847416049401)
		assert.Equal(t, motor.Command{M: motor.RotateLeft, S: motor.Medium}, <-drive) // correcting heading

		heading <- s1.Degree * 168 // heading now corrected to match bearing
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		// waypoint 1
		position <- s2.LatLngFromDegrees(41.185493329876394, -104.80846343169397)
		assert.Equal(t, motor.Command{M: motor.RotateLeft, S: motor.Medium}, <-drive) // correcting heading

		heading <- s1.Degree * 86 // heading now corrected to match bearing
		assert.Equal(t, motor.Command{M: motor.Forward, S: motor.Fast}, <-drive)

		// midpoint 2
		position <- s2.LatLngFromDegrees(41.18550645031859, -104.80825556017193)

		// No heading correction needed, pretty straight shot

		// waypoint 2
		position <- s2.LatLngFromDegrees(41.18551556476105, -104.80804562107231)
		assert.Equal(t, motor.Command{M: motor.Park, S: motor.Stop}, <-drive) // mission complete
	})
}
