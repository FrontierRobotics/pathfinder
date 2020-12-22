#include <Wire.h>
#include <RingBufCPP.h>
#include "Motor.h"

#define I2C_ADDRESS 0x1A
#define I2C_BUFFER_SIZE 32
#define MAX_NUM_ELEMENTS 30
#define MOTOR1_ADDRESS 0x01
#define MOTOR2_ADDRESS 0x02

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
RingBufCPP<struct Event, MAX_NUM_ELEMENTS> buf;

void setup()
{
  motor1.begin();
  motor2.begin();
  Wire.begin(I2C_ADDRESS);
  Wire.onReceive(receiveEvent);
  Wire.onRequest(requestEvent);
}

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

  delay(2000);
}

void requestEvent()
{
  byte sensor_left = analogRead(A0);
  byte sensor_front = analogRead(A1);
  byte sensor_right = analogRead(A2);
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
