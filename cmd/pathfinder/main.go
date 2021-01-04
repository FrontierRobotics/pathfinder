package main

import (
	"bufio"
	"context"
	"image"
	"image/color"
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
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/jacobsa/go-serial/serial"
	"github.com/pbnjay/pixfont"
	"golang.org/x/sync/errgroup"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/devices/ssd1306/image1bit"
	"periph.io/x/periph/host"
)

func C(closer func()error) {
	err := closer()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatalf("host.Init: %v", err)
	}

	serialPort := openUART("/dev/ttyS0")
	defer C(serialPort.Close)

	bus1 := openI2C("1")
	defer C(bus1.Close)

	bus3 := openI2C("3")
	defer C(bus3.Close)

	dev, err := ssd1306.NewI2C(bus1, &ssd1306.Opts{H: 32, W: 128, Sequential: true})
	if err != nil {
		log.Fatalf("failed to initialize ssd1306: %v", err)
	}
	defer C(dev.Halt)

	// see: https://learn.adafruit.com/adafruit-pioled-128x32-mini-oled-for-raspberry-pi/usage
	// see: https://pkg.go.dev/periph.io/x/periph@v3.6.7+incompatible/devices/ssd1306
	// see: https://pkg.go.dev/github.com/pbnjay/pixfont
	// see: https://github.com/pbnjay/pixfont
	// see: https://stackoverflow.com/questions/38299930/how-to-add-a-simple-text-label-to-an-image-in-go
	img := image1bit.NewVerticalLSB(dev.Bounds())
	pixfont.DrawString(img, 0, 0, "L:|||", color.White)
	pixfont.DrawString(img, 0, 8, "F: |", color.White)
	pixfont.DrawString(img, 0, 16, "R: ||", color.White)
	pixfont.DrawString(img, 0, 24, "Test 4 y", color.White)
	if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
		log.Fatal(err)
	}

	var (
		bCtx, cancel = context.WithCancel(context.Background())
		g, ctx       = errgroup.WithContext(bCtx)
		arduino      = &i2c.Dev{Addr: 0x1A, Bus: bus1}
		irReader     = &ir.Reader{
			IRArray: &ir.SensorArray{
				Left:    ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0xA0},
				Forward: ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0xA0},
				Right:   ir.Sensor{ClearUpperBound: 0x10, FarUpperBound: 0xA0},
			},
			Addr: 0x10,
			Tx:   arduino.Tx}
		driver = &motor.Driver{
			Left:  &motor.Motor{Addr: 0x01, Slow: 0x50, Med: 0xA0, Fast: 0xC0},
			Right: &motor.Motor{Addr: 0x02, Slow: 0x50, Med: 0xA0, Fast: 0xC0},
			Tx:    arduino.Tx}
		bno055                     = openBno055(&i2c.Dev{Addr: 0x28, Bus: bus3})
		driverCh                   = make(chan motor.Command, 100)
		irCh                       = make(chan ir.Reading, 10)
		gpsFixCh                   = make(chan bool, 10)
		latLngCh                   = make(chan s2.LatLng, 10)
		headingCh, rollCh, pitchCh = make(chan s1.Angle, 10), make(chan s1.Angle, 10), make(chan s1.Angle, 10)
		pathfinder                 = path.Finder{
			Done:    ctx.Done(),
			GPSfix:  gpsFixCh,
			LatLng:  latLngCh,
			Heading: headingCh,
			Roll:    rollCh,
			Pitch:   pitchCh,
			IR:      irCh,
			Drive:   driverCh,
		}
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
			reader       = bufio.NewReader(serialPort)
			scanner      = bufio.NewScanner(reader)
			lastPosition s2.LatLng
			lastFix      bool
		)
		defer func() { log.Println("ttyS0 loop done") }()
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return nil
			default:
				sentence := scanner.Text()
				r, err := gps.FromGPRMC(sentence)
				if err != nil {
					return err
				}
				if !r.IsEmpty() {
					if lastFix != r.Fix {
						lastFix = r.Fix
						gpsFixCh <- r.Fix
					}
					if lastPosition != r.Position {
						lastPosition = r.Position
						latLngCh <- r.Position
					}
				}
			}
		}
		return nil
	})

	// Start the routine for I2C Bus 1 communication
	g.Go(func() error {
		var (
			irTick        = time.NewTicker(100 * time.Millisecond)
			lastIRReading ir.Reading
		)
		defer func() {
			irTick.Stop()
			// Send a command to park so we don't drive off a cliff
			_ = driver.D(motor.Command{M: motor.Park})
			log.Println("i2c bus 1 loop done")
		}()
		for {
			select {
			case <-ctx.Done():
				return nil
			case cmd := <-driverCh:
				if err := driver.D(cmd); err != nil {
					return err
				}
			case <-irTick.C:
				irReading, err := irReader.Get()
				if err != nil {
					return err
				}
				if irReading != lastIRReading {
					irCh <- irReading
				}
				lastIRReading = irReading
			}
		}
	})

	// Start the routine for I2C Bus 3 communication
	g.Go(func() error {
		var (
			eulerTick                        = time.NewTicker(100 * time.Millisecond)
			bno055TempTick                   = time.NewTicker(10 * time.Second)
			lastHeading, lastRoll, lastPitch s1.Angle
		)
		defer func() {
			eulerTick.Stop()
			bno055TempTick.Stop()
			log.Println("i2c bus 3 loop done")
		}()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-eulerTick.C:
				v, err := bno055.Euler()
				if err != nil {
					return err
				}
				heading, roll, pitch := orientation(v)
				if !heading.ApproxEqual(lastHeading) {
					lastHeading = heading
					headingCh <- heading
				}
				if !roll.ApproxEqual(lastRoll) {
					lastRoll = roll
					rollCh <- roll
				}
				if !pitch.ApproxEqual(lastPitch) {
					lastPitch = pitch
					pitchCh <- pitch
				}
			case <-bno055TempTick.C:
				intTemp, err := bno055.Temperature()
				if err != nil {
					return err
				}
				log.Printf("BNO-055 Temperature: %v °C\n", intTemp)
			}
		}
	})

	g.Go(func() error {
		stopCh := make(chan os.Signal, 1)
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

func openUART(name string) io.ReadWriteCloser {
	serialPort, err := serial.Open(serial.OpenOptions{
		PortName:        name,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	})
	if err != nil {
		log.Fatalf("serial.Open, port: %s err: %v", name, err)
	}
	return serialPort
}

func openI2C(name string) i2c.BusCloser {
	bus, err := i2creg.Open(name)
	if err != nil {
		log.Fatalf("i2c.Open, bus: %s, err: %v", name, err)
	}
	return bus
}
