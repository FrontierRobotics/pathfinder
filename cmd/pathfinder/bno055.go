package main

import (
	"log"
	"time"

	"github.com/andycondon/bno055"
	"github.com/golang/geo/s1"
	"periph.io/x/periph/conn/i2c"
)

type I2CAdapter struct {
	*i2c.Dev
}

func (i *I2CAdapter) Read(reg byte) (byte, error) {
	read := make([]byte, 1)
	err := i.Tx([]byte{reg}, read)
	return read[0], err
}

func (i *I2CAdapter) Write(reg byte, val byte) error {
	return i.Tx([]byte{reg, val}, nil)
}

func (i *I2CAdapter) ReadBuffer(reg byte, buff []byte) error {
	return i.Tx([]byte{reg}, buff)
}

func (i *I2CAdapter) WriteBuffer(reg byte, buff []byte) error {
	write := append([]byte{reg}, buff...)
	return i.Tx(write, nil)
}

func (i *I2CAdapter) Close() error {
	return nil
}

func openBno055(bus *i2c.Dev) *bno055.Sensor {
	log.Printf("*** BNO-055 Startup\n")
	orientationSensor, err := bno055.NewSensorFromBus(&I2CAdapter{
		Dev: bus,
	})
	if err != nil {
		log.Fatalf("bno055.NewSensorFromBus, err: %v", err)
	}

	err = orientationSensor.UseExternalCrystal(true)
	if err != nil {
		log.Fatalf("bno055.UseExternalCrystal, err: %v", err)
	}

	start := time.Now()
	for {
		time.Sleep(1 * time.Second)

		senStatus, err := orientationSensor.Status()
		if err != nil {
			log.Fatalf("bno055.Status, err: %v", err)
		}
		log.Printf("*** Status: system=%v, system_error=%v, self_test=%v\n", senStatus.System, senStatus.SystemError, senStatus.SelfTest)

		if senStatus.SystemError == 0 {
			break
		} else if time.Since(start) > 10*time.Second {
			log.Fatalf("bno055.Status time out")
		}
	}

	revision, err := orientationSensor.Revision()
	if err != nil {
		log.Fatalf("bno055.Revision, err: %v", err)
	}

	log.Printf(
		"*** Revision: software=%v, bootloader=%v, accelerometer=%v, gyroscope=%v, magnetometer=%v\n",
		revision.Software,
		revision.Bootloader,
		revision.Accelerometer,
		revision.Gyroscope,
		revision.Magnetometer,
	)

	axisConfig, err := orientationSensor.AxisConfig()
	if err != nil {
		log.Fatalf("bno055.AxisConfig, err: %v", err)
	}

	log.Printf(
		"*** Axis: x=%v, y=%v, z=%v, sign_x=%v, sign_y=%v, sign_z=%v\n",
		axisConfig.X,
		axisConfig.Y,
		axisConfig.Z,
		axisConfig.SignX,
		axisConfig.SignY,
		axisConfig.SignZ,
	)
	return orientationSensor
}

// See section 3.6.5.4 of the BNO-055 Datasheet
// https://www.bosch-sensortec.com/media/boschsensortec/downloads/datasheets/bst-bno055-ds000.pdf
// x = heading, y = roll, z = pitch
// ex: 2021/01/01 19:06:45 Euler angles: x=0.625, y=-1.312, z=-0.812
//
// Other links:
// see: https://learn.adafruit.com/adafruit-bno055-absolute-orientation-sensor/arduino-code
// see: https://www.bosch-sensortec.com/media/boschsensortec/downloads/application_notes_1/bst-bno055-an007.pdf
// see: https://forum.arduino.cc/index.php?topic=566363.0
// see: https://forums.adafruit.com/viewtopic.php?f=19&t=73014
// see: https://forums.adafruit.com/viewtopic.php?f=19&t=75497&start=15
func orientation(v *bno055.Vector) (s1.Angle, s1.Angle, s1.Angle) {
	var (
		heading = s1.Angle(v.X) * s1.Degree
		roll    = s1.Angle(v.Y) * s1.Degree
		pitch   = s1.Angle(v.Z) * s1.Degree
	)
	return heading, roll, pitch
}
