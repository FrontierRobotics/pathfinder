package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andycondon/pathfinder/pkg/gps"
	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/motor"
	"github.com/andycondon/pathfinder/pkg/path"
	"github.com/andycondon/pathfinder/pkg/status"
	"github.com/jacobsa/go-serial/serial"
	"golang.org/x/sync/errgroup"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func Close(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatalf("host.Init: %v", err)
	}

	// Open /dev/ttyS0 UART serial port
	serialPort, err := serial.Open(serial.OpenOptions{
		PortName:        "/dev/ttyS0",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	})
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer Close(serialPort)

	// Open i2c bus #1
	bus, err := i2creg.Open("1")
	if err != nil {
		log.Fatalf("i2c.Open: %v", err)
	}
	defer Close(bus)

	var (
		bCtx, cancel = context.WithCancel(context.Background())
		g, ctx       = errgroup.WithContext(bCtx)
		arduino      = &i2c.Dev{Addr: 0x1A, Bus: bus}
		irArray      = &ir.SensorArray{
			Left:    ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0xA0},
			Forward: ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0xA0},
			Right:   ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0xA0},
		}
		statusReader = &status.Reader{Addr: 0x10, Tx: arduino.Tx, IRArray: irArray}
		m1           = &motor.Motor{Addr: 0x01, Slow: 0x50, Med: 0xA0, Fast: 0xC0}
		m2           = &motor.Motor{Addr: 0x02, Slow: 0x50, Med: 0xA0, Fast: 0xC0}
		driver       = &motor.Driver{Left: m1, Right: m2, Tx: arduino.Tx, ReadStatus: statusReader.ReadStatus}
		driverCh     = make(chan motor.Command, 100)
		irCh         = make(chan ir.Reading, 100)
		gpsCh        = make(chan gps.Reading, 100)
		stopCh       = make(chan os.Signal, 1)
		pathfinder   = path.Finder{Done: ctx.Done(), GPS: gpsCh, IR: irCh, Drive: driverCh}
	)

	// This is where the magic happens
	g.Go(func() error {
		defer func() { log.Println("pathfinder loop done") }()
		pathfinder.Find()
		return nil
	})

	// Start the routine for reading from the ttyS0 serial port
	g.Go(func() error {
		var (
			reader  = bufio.NewReader(serialPort)
			scanner = bufio.NewScanner(reader)
		)
		defer func() { log.Println("ttyS0 loop done") }()
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return nil
			default:
				sentence := scanner.Text()
				reading, err := gps.FromGPRMC(sentence)
				if err != nil {
					return err
				}
				if !reading.IsEmpty() {
					gpsCh <- reading
				}
			}
		}
		return nil
	})

	// Start the routine for I2C communication
	// Keeps all I2C communication single-threaded
	g.Go(func() error {
		var (
			ticker      = time.NewTicker(10 * time.Millisecond)
			lastReading status.Reading
		)
		defer func() {
			log.Println("i2c loop done")
			// Send a command to park so we don't drive off a cliff
			_, err = driver.D(motor.Command{M: motor.Park})
			ticker.Stop()
		}()
		for {
			select {
			case <-ctx.Done():
				return nil
			case cmd := <-driverCh:
				if _, err := driver.D(cmd); err != nil {
					return err
				}
			case <-ticker.C:
				reading, err := statusReader.Get()
				if err != nil {
					return err
				}

				// TODO Add other I2C sensor checks here

				// Check individual sensors for differences, sending readings on respective channels
				if reading.IR != lastReading.IR {
					irCh <- reading.IR
				}
				lastReading = reading
			}
		}
	})

	g.Go(func() error {
		signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-ctx.Done():
			break
		case <-stopCh:
			log.Println("received shutdown signal...")
			cancel()
		}
		return nil
	})

	err = g.Wait()
	if err != nil {
		log.Fatalf("shut down with error: %v", err)
	} else {
		log.Println("shut down complete")
	}
}
