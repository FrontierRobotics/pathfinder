Arduino Build
=============

Wiring
------

**NOTE!** When connecting and referencing pins, use the Arduino pin number! For example, use pin 2 a.k.a D2 in code
instead of physical pin 4.

Compilation
-----------

From the root of this repository, open [Main.ino](../arduino/Main/Main.ino)

Click the Arduino/Verify button at the top right. (Opt+Cmd+R)

The output of the compilation goes into the [build directory](../arduino/Main/build)(ignored in Git). Inside this directory you'll find two `.hex` files, with and without the bootloader. We can use usually the one without the bootloader because it's already in flash.

Send the compiled `.hex` to the Raspberry pi for flashing:
```
$ scp Main.ino.hex pi@192.168.3.147:/home/pi
```
On the Raspberry Pi:

Flash the `.hex` file:
```
sudo avrdude -P gpio -c gpio -p atmega328p -U flash:w:Main.ino.hex
```

That's it for updating the Arduino firmware. From here the Raspberry Pi code can be executed to communicate with the Arduino over I2C.
