#include <Wire.h>
#include "LCD.h"

#define I2C_ADDRESS 0x1A
#define I2C_BUFFER_SIZE 32
#define LCD_TX_PIN 10
#define LCD_ROWS 4
#define LCD_COLUMNS 20
#define LCD_ADDRESS 0x01

LCD lcd = LCD(LCD_TX_PIN, LCD_ROWS, LCD_COLUMNS);
char i2c_buffer[I2C_BUFFER_SIZE];

void setup() {
  Wire.begin(I2C_ADDRESS);
  Wire.onReceive(receiveEvent);
  Wire.onRequest(requestEvent);
  lcd.begin();
  lcd.clear_screen();
  lcd.set_brightness(0x77);
  lcd.set_cursor(1, 5);
  lcd.print("Pathfinder");
  lcd.set_cursor(2, 7);
  lcd.print("Online");
}

void loop() {
  delay(100);
}

void requestEvent() {
  Wire.write(0x7E);
}

void receiveEvent(int howMany) {
  if (0 >= Wire.available())
  {
    return;
  }

  byte internal_address = Wire.read();
  int data_size = get_data();

  switch (internal_address) {
    case LCD_ADDRESS:
      lcd.print(i2c_buffer);
      break;
  }
}

int get_data()
{
  int index = 0;

  memset(i2c_buffer, 0x00, sizeof(i2c_buffer));

  while (0 < Wire.available()) {
    if (I2C_BUFFER_SIZE <= index)
    {
      break;
    }

    i2c_buffer[index] = Wire.read();
    index++;
  }

  return index;
}

