I2C Protocol
============

I2C Device Address: `0x1A`

To initiate an action or send a command to an attached device, the first byte
must be the internal address for that action. Each action has a unique command
sequence described in the sections below.

The general form for sending commands follows:

| Byte 0               | Byte 1  | Following Bytes... |
| -------------------- | ------- | ------------------ |
| Internal Address     | Command | Command Parameters |

The available actions and their internal address:

| Action  | Internal Addr |
| ------- | ------------- |
| LCD     | `0x00`        |
| Motor 1 | `0x01`        |
| Motor 2 | `0x02`        |
| Status  | `0x10`        |

LCD
---

**_Not Implemented_** - not necessary, and somewhat unsafe to transmit arbitrary message lengths.

The second byte in the transmission (Byte 1) specifies the command for the LCD.

| Command        | Byte 1 | Parameters                |
| -------------- | ------ | ------------------------- |
| Write          | `0x00` | 20 or less bytes to print |
| Set Cursor     | `0x01` | byte: row, byte: column   |
| Set Brightness | `0x02` | byte: level               |
| Clear          | `0x02` | None                      |

Motor
-----

The second byte in the transmission (Byte 1) specifies the command for the
Motor.

| Command | Byte 1 | Parameters                   |
| ------- | ------ | ---------------------------- |
| Forward | `0x00` | byte: speed, `0x00` for stop |
| Reverse | `0x01` | byte: speed, `0x00` for stop |

Status
------

A status message is sent to the master with every interaction. It is possible
to receive a dedicated status message as well.

The form of the status message is a 1 byte value for each sensor.

|        | Byte 0 | Byte 1 | Byte 2 |
| ------ | ------ | ------ | ------ |
| Sensor | IR 1   | IR 2   | IR 3   |
