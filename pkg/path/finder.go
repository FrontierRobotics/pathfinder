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
	Done  <-chan struct{}
}

func (f *Finder) Find() {
	var (
		IR      ir.Reading
		GPS     gps.Reading
		lastCmd motor.Command
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
