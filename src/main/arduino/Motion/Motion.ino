#include <Wire.h>
#include "LCD.h"
#include "Motor.h"

#define I2C_ADDRESS 0x1A
#define I2C_BUFFER_SIZE 32
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

// Remember to use Arduino pins, not physical ones
Motor motor1 = Motor(2, 3);
Motor motor2 = Motor(4, 5);

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

void loop()
{
  delay(100);
}

void requestEvent()
{
  Wire.write("howdy");
}

void receiveEvent(int receive_size)
{
  if (!Wire.available())
  {
    return;
  }
  byte internal_address = Wire.read();

  switch (internal_address)
  {
  case LCD_ADDRESS:
    lcd_event(receive_size - 1);
    break;
  case MOTOR1_ADDRESS:
    // Display Reference
    // 0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19
    // M  1  :     F  0  5  5     M  2  :     R  0  5  5
    motor_event(&motor1, receive_size - 1);
    lcd.set_cursor(1, 0);
    lcd.print("M1: %s", motor1.status());
    break;
  case MOTOR2_ADDRESS:
    motor_event(&motor2, receive_size - 1);
    lcd.set_cursor(1, 9);
    lcd.print("M2: %s", motor2.status());
    break;
  }
}

void lcd_event(int receive_size)
{
  if (!Wire.available())
  {
    return;
  }
  byte command = Wire.read();

  switch (command)
  {
  case LCD_COMMAND_WRITE:
  {
    int data_size = get_data();
    if (data_size > 0)
    {
      lcd.print(i2c_buffer);
    }
    break;
  }
  case LCD_COMMAND_SET_CURSOR:
  {
    if (!Wire.available())
    {
      return;
    }
    byte row = Wire.read();
    if (!Wire.available())
    {
      return;
    }
    byte column = Wire.read();
    lcd.set_cursor(row, column);
    break;
  }
  case LCD_COMMAND_SET_BRIGHTNESS:
  {
    if (!Wire.available())
    {
      return;
    }
    byte brightness = Wire.read();
    lcd.set_brightness(brightness);
    break;
  }
  case LCD_COMMAND_CLEAR:
  {
    lcd.clear_screen();
    lcd.set_cursor(0, 0);
    break;
  }
  }
}

void motor_event(Motor *motor, int receive_size)
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
  if (0x01 == direction)
  {
    motor->forward(speed);
  }
  else
  {
    motor->reverse(speed);
  }
}

int get_data()
{
  int index = 0;

  memset(i2c_buffer, 0x00, sizeof(i2c_buffer));

  while ((0 < Wire.available()) && (index <= I2C_BUFFER_SIZE))
  {
    i2c_buffer[index] = Wire.read();
    index++;
  }

  return index;
}
