package main

import (
	"log"
	"time"

	"github.com/andycondon/pathfinder/pkg/path"
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
	defer bus.Close()

	a := NewArduino(bus)

	for {
		status, err := a.GetStatus()
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}

		p := path.Find(status.IR)

		log.Printf("%s Path:%v\n", status.IR.String(), p)
		time.Sleep(time.Millisecond * 1000)
	}
}