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
  lcd.clear_screen();
  while (1 < Wire.available()) { // loop through all but the last
    char c = Wire.read(); // receive byte as a character
    lcd.print("%c", c);
  }
  int x = Wire.read();    // receive byte as an integer

  lcd.print("%d", x);
}
