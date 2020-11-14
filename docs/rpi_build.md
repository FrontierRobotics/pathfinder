Raspberry Pi Build
==================

Cross Compiling
---------------

Change directories to where `main.go` resides and run `go build` with environment variables for the target platform:
```
$ cd cmd/pathfinder
$ env GOOS=linux GOARCH=arm GOARM=5 go build
```

Send the compiled binary to the Raspberry Pi:
```
$ scp pathfinder pi@192.168.3.147:/home/pi
```

Executing
---------

First make sure the Arduino code is [built and flashed](arduino_build.md).

On the Raspberry Pi:

Run the executable to start upPathfinder:

```
$ ./pathfinder
```

Or, to test the I2C commands directly...

Make both motors go forward at 50%:
```
$ ./test-i2c 010180 // Motor 1
$ ./test-i2c 020180 // Motor 2
```

