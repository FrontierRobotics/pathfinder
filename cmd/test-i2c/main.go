package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

// Flash an LED
//func main() {
//	host.Init()
//	t := time.NewTicker(500 * time.Millisecond)
//	for l := gpio.Low; ; l = !l {
//		rpi.P1_33.Out(l)
//		<-t.C
//	}
//}

func main() {
	args := os.Args[1:]
	if 0 == len(args) {
		log.Fatal("please pass a hex encoded string")
	}
	write, err := hex.DecodeString(args[0])
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("write: ")
	for _, b := range write {
		fmt.Printf("0x%02X ", b)
	}
	log.Printf("\n")

	if _, err := host.Init(); err != nil {
		log.Fatalf("%v", err)
	}

	bus, err := i2creg.Open("1")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer bus.Close()

	dev := &i2c.Dev{Addr: 0x1a, Bus: bus}
	read := make([]byte, 3)

	if err := dev.Tx(write, read); err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("read: %v\n", read)
}
