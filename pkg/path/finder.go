package path

import (
	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/motor"
	"github.com/andycondon/pathfinder/pkg/nav"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"log"
)

var (
	fast    = motor.Command{M: motor.Forward, S: motor.Fast}
	forward = motor.Command{M: motor.Forward, S: motor.Slow}
	left    = motor.Command{M: motor.RotateLeft, S: motor.Medium}
	right   = motor.Command{M: motor.RotateRight, S: motor.Medium}
	park    = motor.Command{M: motor.Park}
)

// Helpful links:
// see: https://diydrones.com/profiles/blogs/the-difference-between-heading
// see: https://www.movable-type.co.uk/scripts/latlong.html
type Finder struct {
	IR                   <-chan ir.Reading
	GPSfix               <-chan bool
	Position             <-chan s2.LatLng
	Heading, Roll, Pitch <-chan s1.Angle
	Drive                chan<- motor.Command
	Done                 <-chan struct{}
	Waypoints            []s2.LatLng
	BearingTolerance     s1.Angle
}

func (f *Finder) Find() {
	var (
		IR                            ir.Reading
		gpsFix                        bool
		position                      s2.LatLng
		heading, roll, pitch, bearing s1.Angle
		distance                      nav.Distance
		lastCmd                       motor.Command
		mission                       = NewMission(f.Waypoints)
		drive                         = func(cmd motor.Command) {
			if lastCmd != cmd {
				lastCmd = cmd
				f.Drive <- cmd
			}
		}
	)
	for {
		waypoint := mission.CurrentWaypoint()
		select {
		case <-f.Done:
			return
		case gpsFix = <-f.GPSfix:
			if gpsFix {
				log.Println("GPS fix acquired ")
			} else {
				log.Println("GPS fix lost ")
			}
		case position = <-f.Position:
			waypoint = mission.Progress(position)
			bearing = nav.InitialBearing(position, waypoint.LatLng)
			distance = nav.DistanceBetweenPositions(position, waypoint.LatLng)
			log.Printf("position: %s, heading: %s, bearing: %s, distance: %s\n", position.String(), heading.String(), bearing.String(), distance.String())
		case heading = <-f.Heading:
			log.Printf("heading: %s°\n", heading.String())
		case roll = <-f.Roll:
			log.Printf("roll: %s°\n", roll.String())
		case pitch = <-f.Pitch:
			log.Printf("pitch: %s°\n", pitch.String())
		case IR = <-f.IR:
			log.Printf("IR: %s", IR.String())
		}

		if !gpsFix || mission.Complete() {
			drive(park)
			continue
		}

		obstacle, cmd := avoid(IR)
		relBearing := nav.RelativeBearing(heading, bearing, f.BearingTolerance)
		if !obstacle && relBearing != 0 {
			log.Printf("rel heading: %s°\n", heading.String())
			if relBearing > 0 {
				cmd = left
			} else {
				cmd = right
			}
		}

		drive(cmd)
	}
}

func avoid(IR ir.Reading) (bool, motor.Command) {
	if IR.AllClear() {
		return false, fast
	}
	if !IR.F.IsNear() {
		if !IR.R.IsNear() && !IR.L.IsNear() {
			return false, forward
		}
		if !IR.R.IsNear() {
			return true, right
		}
		return true, left
	}
	if !IR.L.IsNear() {
		return true, left
	}
	return true, right
}
