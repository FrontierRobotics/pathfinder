package main

import (
	"log"
	"time"

	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/path"
	"github.com/andycondon/pathfinder/pkg/status"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatalf("%v", err)
	}

	// Open i2c bus #1
	bus, err := i2creg.Open("1")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer func() {
		if err := bus.Close(); err != nil {
			log.Fatalf("%v", err)
		}
	}()

	var (
		arduino = &i2c.Dev{Addr: 0x1A, Bus: bus}
		irArray = &ir.SensorArray{
			Left:    ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0x50},
			Forward: ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0x50},
			Right:   ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0x50},
		}
		s      = &status.Reader{Address: 0x10, Tx: arduino.Tx, IRArray: irArray}
	)

	var lastReading status.Reading
	for {
		reading, err := s.GetStatus()
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}

		// Arduino status won't be the only input for path finding, so always find path based on all inputs.
		p := path.Find(reading.IR)

		if reading != lastReading {
			log.Printf("%s Path:%v\n", reading.IR.String(), p)
		}
		lastReading = reading


		time.Sleep(time.Millisecond * 1000)
	}
}
