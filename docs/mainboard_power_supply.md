# Mainboard Power Supply

The power supply for the mainboard is designed with the following features:

* Wide voltage input range of 5.4V to 36V.
* Selectively disabled 2A, 5V and 1A, 3.3V power buses.
* Always-On 1A, 5V critical power bus.

## Circuit Overview

![Circuit Diagram](images/mainboard_power_supply.png)

The heart of the power supply is the LT3514 Triple Step-Down Switching Regulator. It consists of three buck regulators with a single 2A channel, and two 1A channels. The device is capable of accepting voltages from 3.2V up to 36V. An on-chip boost regulator allows each channel to operate up to 100% duty cycle, minimizing the need for leveling circuitry.

The LT3514 operates robustly in fault conditions. Cycle-by-cycle peak current limit and catch diode current limit sensing protect the IC, and other mainboard components during overload conditions. Thermal shutdown protects the internal power switches at elevated temperatures. Soft-start helps control the peak inductor current during startup.

The LT3514 features output voltage tracking and sequencing, programmable frequency, programmable undervoltage lockout, and a power good pin to indicate when all ouputs are in regulation.

Channel 3 has the highest current rating, and thus will be used to supply the Raspberry Pi and 5V bus. Channel 1, at 1A, will supply the 3.3V bus. Channel 4 will supply the 5V critical bus for the mainboard power monitor. Channels 1 and 3 will be selectively disabled by the mainboard power monitor circuit during shutdown or low power scenarios. Channel 4 will always be on, supplying the power monitor in standby. By that nature, EN/UVLO will always be enabled as well to ensure constant operation.

## FB Resistor Network

The general formula for the FB resistor network is:

```
R1 = R2 * ((Vout/0.8V) - 1)
```

The datasheet states to use 1% tolerance resisitors, and that a good value for R2 is 10.2kΩ. R2 should not exceed 20kΩ to avoid bias current error.

therefore:

```
Channel 1: 1A, 3.3V
R1 = 10.2kΩ * ((3.3V/0.8V - 1) = 31875Ω ≅ 31.6kΩ

Channel 3: 2A, 5V
R1 = 10.2kΩ * ((5V/0.8V - 1) = 53550Ω ≅ 53.6kΩ

Channel 4: 1A, 5V
R1 = 10.2kΩ * [(5V/0.8V) - 1] = 53550Ω ≅ 53.6kΩ
```
conclusion:

```
R16,18,20 = 10.2kΩ @ ±1%
R15       = 31.6kΩ @ ±1%
R17,19    = 53.6kΩ @ ±1%
```

## Input Voltage Range

The minimum input voltage to regulate the output generally has to be at least 400mV greater than the greatest programmed output voltage. The only exception is when the largest output is less than 2.8V in which case the minimum is 3.2V

The absolute maximum input voltage is 40V, and the LT3514 will regulate so long as the voltage remains less than or equal to that value. However, for constant frequency operation the maximum input voltage is determined through the following formula:

```
VIN(PS) = [(VOUT + VD)/(fSW * tON(MIN)] + VSW - VD
```

where:

* `VIN(PS)` is the maximum input voltage to operate in constant frequency operation without skipping pulses.
* `VOUT` is the programmed output voltage.
* `VSW` is the switch voltage drop, at `IOUT1,4 = 1A`, `VSW1,4 = 0.4V`, at `IOUT3 = 2A`, `VSW3 = 0.4V`.
* `VD` is the catch diode forward voltage drop, for an appropriately sized diode, `VD = 0.4V`.
* `fSW` is the programmed switching frequency. As will be shown in the next section, `fSW = 1.0 MHz`.
* `tON(MIN)` is the minimum on-time, worst-case over temperature = 110ns (at T = 125ºC).

therefore:

```
VIN(PS) = [(5V + 0.4V)/(1.0 MHz * 110ns)] + 0.4V - 0.4V = 49.091V = 40V (device maximum)
VIN(PS) = [(3.3V + 0.4V)/(1.0 MHz * 110ns)] + 0.4V - 0.4V = 40.7V = 40V (device maximum)
```

conclusion:

Our computed `VIN(PS) = 40.7V` exceeds the maximum of 40V, and thus 40V shall be the upper value of VIN. However, the datasheet recommends not starting the LT3514 at input voltages greater than 36V, as the LT3514 must simultaneously conduct maximum currents at high VIN.

Therefore, the recommended input voltage range is `5.4V < VIN < 36V`.

## Frequency Selection

There are two ways to program the frequency:

* Tying a 1% tolerance resistor `RT` from the `RT/SYNC` pin to ground.
* Synchronize the internal oscillator to an external clock.

For our power supply, we'll use the first method. The 1.0 MHz frequency choosen is in the midrange of the LT3514's 250kHz to 2.5MHz operating range. In addition, as shown above, 1.0 MHz results in a full-range `VIN(PS)`. The datasheet includes a table for choosing the proper resistor value for different frequencies.

conclusion:

```
fSW = 1.0 MHz
RT = 18.2kΩ @ ±1%
```

## BOOST Regulator and SKY Pin Considerations

The on-chip boost regulator generates the SKY voltage to be 4.85V above VIN. The SKY voltage is the source of the drive current for the buck regulators that drive the output channels. A good choice for the inductor that will ensure each buck regulator will have sufficient drive current is given by:

```
L = 20.5µH / f
```

where `f` is in MHz. This gives a value of `L = 20.5µH` for `f = 1.0 MHz`.

The optimal SKY pin output current requirement and inductor value is calculated by:

```
ISKY = (IOUT1 * VOUT1 + IOUT3 * VOUT3 + IOUT4 * VOUT4) / 50 * VIN

L = VIN * DC5 / 2 * fSW * [0.3 * (1 - 0.25 * DC5) - ISKY]
```

where:

* `fSW = 1.0 MHz` is the programmed switching frequency.
* `VIN = 12V`, which will be a typical value.
* `DC5 = 0.29412` is the boost regulator duty cycle, given by: `DC5 = 5V/(VIN + 5V)`.
* `VOUT3,4 = 5V`, `VOUT1 = 3.3V`, `IOUT1,4 = 1A`, `IOUT3 = 2A`

therefore:

```
ISKY = (1A * 3.3V + 2A * 5V + 1A * 5V) / 50 * 12
ISKY = 18.3 / 600 = 30.5mA

L = 12V * 0.29412 / 2 * 1.0 MHz * [0.3 * (1 - 0.25 * 0.29412) - 30.5mA]
L = 3.52944 / 2000000 * (0.3 * 0.92647 - 0.0305)
L = 3.52944 / 2000000 * 0.247441 = 0.00000713188194H
L ≅ 7.1µH
```

From the example in the datasheet using different output voltages, we calculated this optimal value of L:

```
ISKY = 13.4 / 600 = 22.3mA
L ≅ 6.9µH
```

This is very close to our computed value. The value of L chosen for the example is 10µH. We will also use this value for our design.

conclusion:

The chosen value of the inductor is 10µH for `VIN = 12V`.

## Inductor Selection

```
L = 2 * (VOUT + VD)/fSW for Channels 1, 4
L = (VOUT + VD)/fSW for Channel 3
```

## Bill of Materials
