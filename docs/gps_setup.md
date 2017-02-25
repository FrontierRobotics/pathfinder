# GPS Setup

Pathfinder uses the [Ultimate GPS module from Adafruit](https://www.adafruit.com/product/746).

## Stealing back the serial port from the console

Aside from enabling Serial like we did in [Raspberry Pi Setup](raspberry_pi_setup.md), we have to take it away from the console service.

```
sudo systemctl stop serial-getty@ttyS0.service
sudo systemctl disable serial-getty@ttyS0.service
```

You also need to remove the console from the cmdline.txt. If you edit this with:

```
sudo nano /boot/cmdline.txt
```

You will see something like:

```
dwc_otg.lpm_enable=0 console=serial0,115200 console=tty1 root=/dev/mmcblk0p2 rootfstype=ext4 elevator=deadline fsck.repair=yes root wait
```

Remove the part mapping the console to serial0:

```
console=serial0,115200
```

Save and reboot for changes to take effect.

## Observing the output

The Ultimate GPS module runs at 9600 baud 8N1.

### Minicom

If this is your thing, first start it up in setup mode:

```
sudo minicom -s
```

Put in the proper device, `/dev/ttyS0` and the baud rate, then start it up like normal:

```
sudo minicom
```

### Screen

This is probably the easiest.

```
sudo screen /dev/ttyS0 9600
```

That's it! [See the rest of the documentation to continue.](../README.md)

----

* [Using UART](https://learn.adafruit.com/adafruit-ultimate-gps-on-the-raspberry-pi/using-uart-instead-of-usb)
* [Configuring GPIO Serial](http://spellfoundry.com/2016/05/29/configuring-gpio-serial-port-raspbian-jessie-including-pi-3)
