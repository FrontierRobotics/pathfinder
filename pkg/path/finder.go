package path

import (
	"log"

	"github.com/andycondon/pathfinder/pkg/gps"
	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/motor"
)

type Finder struct {
	IR    <-chan ir.Reading
	GPS   <-chan gps.Reading
	Drive chan<- motor.Command
	Done  chan struct{}
}

func (f *Finder) Find() {
	var (
		IR  ir.Reading
		GPS gps.Reading
	)
	for {
		select {
		case <-f.Done:
			return
		case GPS = <-f.GPS:
			log.Println("finder gps - " + GPS.String())
		case IR = <-f.IR:
			log.Println("finder ir - " + IR.String())
		}

		f.Drive <- avoid(IR)
	}
}

func avoid(IR ir.Reading) motor.Command {
	var (
		forward = motor.Command{M: motor.Forward, S: motor.Slow}
		left    = motor.Command{M: motor.RotateLeft, S: motor.Medium}
		right   = motor.Command{M: motor.RotateRight, S: motor.Medium}
	)
	if !IR.F.IsNear() && !IR.L.IsNear() && !IR.R.IsNear() {
		return forward
	}
	if !IR.F.IsNear() && IR.L.IsNear() && !IR.R.IsNear() {
		return right
	}
	if !IR.F.IsNear() && !IR.L.IsNear() && IR.R.IsNear() {
		return left
	}
	if IR.F.IsNear() && !IR.L.IsNear() {
		return left
	}
	if IR.F.IsNear() && !IR.R.IsNear() {
		return right
	}
	return right
}
