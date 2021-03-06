# Raspberry Pi Setup

## Load the Image

1. Download the latest [Raspbian Jessie with Pixel image](https://downloads.raspberrypi.org/raspbian_latest).
2. Use [Etcher](https://etcher.io/) to burn the image onto the SD card.

Once you've got the image loaded hook up a keyboard, mouse, TV via HDMI, and USB power supply. Then power it up!

## Set up WiFi

### Step 1: Give yourself administrative access:

```
sudo su
```

### Step 3: Identify your network adapter and SSID:

Note: Take note of the line with 802.11 on it and what comes to the left of it. This is your wireless card and you'll need it in the steps that follow.

```
iwconfig
```

The output of this command will look something like the following:

```
eth0 no wireless extensions.

wlan0 IEEE 802.11bgn ESSID:off/any
Mode:Managed Access Point: Not-Associated Tx-Power=20 dBm
Retry short limit:7 RTS thr:off Fragment thr:off
Power Management:off

lo no wireless extensions.
```

Run the following command to make sure your SSID is visible to the interface:

```
sudo iwlist wlan0 scan | grep "ESSID"
```

### Step 4: Stop Network Manager if it is running:

```
service network-manager stop
```

### Step 5: Create a wpa_supplicant.conf file by replacing access_point_name with your SSID:

To avoid storing your password in plain text, run the following command

```
wpa_passphrase "access_point_name" > /etc/wpa_supplicant/wpa_supplicant_temp.conf
```

Note: The prompt you'd normally get is eaten by the redirection to the file. Just type your password and everything will be ok.

Copy the content of what that produced at the bottom of `/etc/wpa_supplicant/wpa_supplicant_temp.conf`. You'll want to remove your commented-out plain text password.

TODO Paste in example

### Step 6: Set up the network interface file:

```
nano /etc/network/interfaces
```

Edit the existing interface, or set up a new one at the bottom of the file (if your wireless adapter does not show up as wlan0 replace it with what did show up when running iwconfig above) to get an IP address automatically via your DHCP/router/access point:

```
auto wlan0
iface wlan0 inet dhcp
wpa-conf /etc/wpa_supplicant/wpa_supplicant.conf
```

### Step 7: Activate the wireless connection:

```
ifdown wlan0 && ifup wlan0
```

Sometimes it fails on the first try. Give it another go before banging your head too hard.

## Configuration

Turn on the appropriate settings depending on what you'll need. Don't turn on stuff you won't use, it'll just take more energy. Once finished click OK. It'll probably want to reboot.

Go into Start->Preferences->Raspberry Pi Configuration

### System

Set Password to: <Your Password> (default is raspberry)

Set Hostname to: pathfinder

### Interfaces

Turn on SSH, SPI, I2C, Serial.

### Performance

Nothing to change.

### Localisation

Set Locale to:

* Language: en (English)
* Country: US (USA)
* Character Set: UTF-8

Set Timezone to:

* Area: Etc
* Location: UTC

Set Keyboard to: United States, English (US)

Set WiFi Country Code to: US United States

## Get the IP and log in remotely

First find your Raspberry Pi's IP address running the following command:

```
ifconfig
```

Or use `nmap` to find it on the network:

```
nmap -sn 192.168.1.0/24

Starting Nmap 7.40 ( https://nmap.org ) at 2017-02-24 17:42 MST
.
.
Nmap scan report for pathfinder (192.168.1.16)
Host is up (0.038s latency).
.
.
Nmap done: 256 IP addresses (9 hosts up) scanned in 2.41 seconds
```

Once you've got the IP, log in through ssh:

```
ssh pi@192.168.1.16
```

## Upgrade and customize the OS

Run `sudo apt-get update` and `sudo apt-get upgrade` to get the system up-to-date.

Install the following packages:

```
sudo apt-get install vim screen minicom
```

Then edit `~/.bashrc` and uncomment the following lines:

```
# some more ls aliases
alias ll='ls -l'
alias la='ls -A'
alias l='ls -CF'
```

That's it! [See the rest of the documentation to continue.](../README.md)

----

### References

* [ThinkPenguin.com - Wireless CLI](https://www.thinkpenguin.com/gnu-linux/how-configure-wifi-card-using-command-line-or-terminal)
* [RaspberryPi.org - Wireless CLI](https://www.raspberrypi.org/documentation/configuration/wireless/wireless-cli.md)
* [RaspberryPi.org - IP Address](https://www.raspberrypi.org/documentation/remote-access/ip-address.md)
