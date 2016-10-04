#include <Wire.h>
#include "LCD.h"

#define I2C_ADDRESS 0x1A
#define LCD_TX_PIN 10
#define LCD_ROWS 4
#define LCD_COLUMNS 20

LCD lcd = LCD(LCD_TX_PIN, LCD_ROWS, LCD_COLUMNS);

void setup() {
  Wire.begin(I2C_ADDRESS);
  Wire.onReceive(receiveEvent);
  Wire.onRequest(requestEvent);
  lcd.begin();
  lcd.clear_screen();
}

void loop() {
  lcd.set_brightness(0x00);
  delay(1000);
  lcd.set_brightness(0x22);
  delay(1000);
  lcd.set_brightness(0x44);
  delay(1000);
  lcd.set_brightness(0xFF);
  delay(1000);

  lcd.set_cursor(1, 5);
  lcd.print("Hello!");
  delay(1000);
  lcd.set_cursor(2, 10);
  lcd.print("Howdy!");
  delay(1000);
  lcd.blink_cursor();
  delay(1000);
  lcd.disable_cursor();
  lcd.clear_screen();
}

void requestEvent() {
  Wire.write(0x7E);
}

void receiveEvent(int howMany) {
  while (1 < Wire.available()) { // loop through all but the last
    char c = Wire.read(); // receive byte as a character
  }
  int x = Wire.read();    // receive byte as an integer
}
