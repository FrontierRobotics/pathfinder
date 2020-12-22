#include <Wire.h>
#include <RingBufCPP.h>
#include "LCD.h"
#include "Motor.h"

#define I2C_ADDRESS 0x1A
#define I2C_BUFFER_SIZE 32
#define MAX_NUM_ELEMENTS 30
#define LCD_TX_PIN 10
#define LCD_ROWS 4
#define LCD_COLUMNS 20
#define LCD_ADDRESS 0x00
#define LCD_COMMAND_WRITE 0x00
#define LCD_COMMAND_SET_CURSOR 0x01
#define LCD_COMMAND_SET_BRIGHTNESS 0x02
#define LCD_COMMAND_CLEAR 0x03
#define MOTOR1_ADDRESS 0x01
#define MOTOR2_ADDRESS 0x02

LCD lcd = LCD(LCD_TX_PIN, LCD_ROWS, LCD_COLUMNS);
char i2c_buffer[I2C_BUFFER_SIZE];

struct Event
{
  byte internal_address;
  byte direction;
  byte speed;
};

// Remember to use Arduino pins, not physical ones
Motor motor1 = Motor(2, 3);
Motor motor2 = Motor(4, 5);
volatile byte sensor_left, sensor_front, sensor_right;
RingBufCPP<struct Event, MAX_NUM_ELEMENTS> buf;

void setup()
{
  motor1.begin();
  motor2.begin();
  Wire.begin(I2C_ADDRESS);
  Wire.onReceive(receiveEvent);
  Wire.onRequest(requestEvent);
  lcd.begin();
  lcd.clear_screen();
  lcd.set_brightness(0x77);
  lcd.set_cursor(0, 5);
  lcd.print("Pathfinder");
}

// Display Reference
// 0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19
//                   P  a  t  h  f  i  n  d  e  r
// M  1  :     F  0  5  5     M  2  :     R  0  5  5
//
// L  :  0  5  5     F  :  0  5  5     R  :  0  5  5
void loop()
{
  struct Event e;

  while (buf.pull(&e))
  {
    switch (e.internal_address)
    {
    case MOTOR1_ADDRESS:
      motor_event(&motor1, e.direction, e.speed);
      break;
    case MOTOR2_ADDRESS:
      motor_event(&motor2, e.direction, e.speed);
      break;
    }
  }
  lcd.set_cursor(1, 0);
  lcd.print("M1: %s", motor1.status());
  lcd.set_cursor(1, 9);
  lcd.print("M2: %s", motor2.status());
  lcd.set_cursor(3, 0);
  lcd.print("L:%03d F:%03d R:%03d", sensor_left, sensor_front, sensor_right);
  delay(2000);
}

void requestEvent()
{
  sensor_left = analogRead(A0);
  sensor_front = analogRead(A1);
  sensor_right = analogRead(A2);
  Wire.write(sensor_left);
  Wire.write(sensor_front);
  Wire.write(sensor_right);
}

void receiveEvent(int receive_size)
{
  if (!Wire.available())
  {
    return;
  }
  byte internal_address = Wire.read();

  if (internal_address == MOTOR1_ADDRESS || internal_address == MOTOR2_ADDRESS)
  {
    if (!Wire.available())
    {
      return;
    }
    byte direction = Wire.read();

    if (!Wire.available())
    {
      return;
    }
    byte speed = Wire.read();

    struct Event e;
    e.internal_address = internal_address;
    e.direction = direction;
    e.speed = speed;
    buf.add(e);
  }
}

void motor_event(Motor *motor, byte direction, byte speed)
{
  if (0x01 == direction)
  {
    motor->forward(speed);
  }
  else
  {
    motor->reverse(speed);
  }
}
