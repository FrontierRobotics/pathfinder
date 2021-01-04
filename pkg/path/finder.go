package path

import (
	"log"

	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/motor"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

type Finder struct {
	IR                   <-chan ir.Reading
	GPSfix               <-chan bool
	LatLng               <-chan s2.LatLng
	Heading, Roll, Pitch <-chan s1.Angle
	Drive                chan<- motor.Command
	Done                 <-chan struct{}
}

func (f *Finder) Find() {
	var (
		IR                   ir.Reading
		gpsFix               bool
		latLng               s2.LatLng
		heading, roll, pitch s1.Angle
		lastCmd              motor.Command
	)
	for {
		select {
		case <-f.Done:
			return
		case gpsFix = <-f.GPSfix:
			if gpsFix {
				log.Println("GPS fix acquired ")
			}
		case latLng = <-f.LatLng:
			log.Printf("finder position - %s\n", latLng.String())
		case heading = <-f.Heading:
			log.Printf("finder heading - %s°\n", heading.String())
		case roll = <-f.Roll:
			log.Printf("finder roll - %s°\n", roll.String())
		case pitch = <-f.Pitch:
			log.Printf("finder pitch - %s°\n", pitch.String())
		case IR = <-f.IR:
			log.Printf("finder ir - %s", IR.String())
		}

		if cmd := avoid(IR); cmd != lastCmd {
			lastCmd = cmd
			f.Drive <- cmd
		}
	}
}

func avoid(IR ir.Reading) motor.Command {
	var (
		fast    = motor.Command{M: motor.Forward, S: motor.Fast}
		forward = motor.Command{M: motor.Forward, S: motor.Slow}
		left    = motor.Command{M: motor.RotateLeft, S: motor.Medium}
		right   = motor.Command{M: motor.RotateRight, S: motor.Medium}
	)
	if IR.AllClear() {
		return fast
	}
	if !IR.F.IsNear() {
		if !IR.R.IsNear() && !IR.L.IsNear() {
			return forward
		}
		if !IR.R.IsNear() {
			return right
		}
		return left
	}
	if !IR.L.IsNear() {
		return left
	}
	return right
}
