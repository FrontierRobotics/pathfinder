#include <Wire.h>

#define I2C_ADDRESS 0x1A

void setup() {
  Wire.begin(I2C_ADDRESS);
  Wire.onReceive(receiveEvent);
  Wire.onRequest(requestEvent);
  lcd_initialize();
  lcd_clear();
}

void loop() {
  lcd_set_brightness(0x00);
  delay(1000);
  lcd_set_brightness(0x22);
  delay(1000);
  lcd_set_brightness(0x44);
  delay(1000);
  lcd_set_brightness(0xFF);
  delay(1000);

  lcd_set_cursor(1, 5);
  lcd_print("Hi!");
  delay(1000);
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
