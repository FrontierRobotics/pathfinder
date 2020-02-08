# Atmel Programmer Setup

Steve Marple designed an implementation of a AVR ISP programmer using the Raspberry Pi GPIO port which can be used to program Atmel's AVR range of microcontrollers with avrdude. An ISP programmer based on this design was incorporated into the Pathfinder mainboard. Mr. Marple usThis post explains how to use avrdude to actually program devices.

## Software installation
Use `scp` to upload the [patched version of avrdude](../avrdude) onto the Pi. You will probably want the armhf (hardware floating-point) version. Upload the documentation package for avrdude too. Install the packages using `dpkg -i`. For example:

```
sudo dpkg -i avrdude_5.10-4_armhf.deb
sudo dpkg -i avrdude-doc_5.10-4_all.deb
```

Using avrdude over the GPIO interface is problematic for users other than root. The easiest solution is to give the avrdude binary setgroup permission:

```
sudo chmod g+s /usr/bin/avrdude
```

## Usage

Selecting the GPIO programmer is simply a matter of including `-P gpio -c gpio` options; the `-P` option specifies that the GPIO port is used (as opposed to USB, serial or parallel interfaces) whilst the `-c` option selects the correct programmer type on that port.

For example, to check the signature on an ATmega328P execute the command:

```
sudo avrdude -P gpio -c gpio -p atmega328p
```

To read the fuses execute the command:

```
sudo avrdude -P gpio -c gpio -p atmega328p -U lfuse:r:-:h -U hfuse:r:-:h -U efuse:r:-:h
```

To flash the device:

```
sudo avrdude -P gpio -c gpio -p atmega328p -U flash:w:Motion.ino.standard.hex
```
To read the device's flash:

```
sudo avrdude -P gpio -c gpio -p atmega328p -U flash:r:backup.hex:r
```

## Customization

The packages above define a single programmer called gpio which uses the gpio interface on GPIO pins 8 to 11. If expanding the mainboard to implement two independent programmers, use gpio0 and gpio1. You can add these by creating a `.avrduderc` file in your home directory. The file should contain:

```
programmer
  id    = "gpio0";
  desc  = "Use sysfs interface to bitbang GPIO lines";
  type  = gpio;
  reset = 8;
  sck   = 11;
  mosi  = 10;
  miso  = 9;
;

programmer
  id    = "gpio1";
  desc  = "Use sysfs interface to bitbang GPIO lines";
  type  = gpio;
  reset = 7;
  sck   = 11;
  mosi  = 10;
  miso  = 9;
;
```

That's it! [See the rest of the documentation to continue.](../README.md)

----

### References

* [Steve Marple - ISP Design](http://blog.stevemarple.co.uk/2012/07/avrarduino-isp-programmer-using.html)
* [Steve Marple - How to Use](http://blog.stevemarple.co.uk/2013/03/how-to-use-gpio-version-of-avrdude-on.html)
* [Steve Marple - RPi_RFM12B_ISPd](https://github.com/stevemarple/RPi_RFM12B_ISP/tree/master/software/avrdude)
* [Sparkfun](https://learn.sparkfun.com/tutorials/pocket-avr-programmer-hookup-guide/using-avrdude)
