#include <Arduino.h>
#include <SoftwareSerial.h>
#include "LCD.h"

#define LCD_RX_PIN 30  // rxPin is immaterial - not used - just make this an unused Arduino pin number

LCD::LCD(int txPin,byte rows, byte columns) : _serial(SoftwareSerial(LCD_RX_PIN, txPin)) {
  pinMode(txPin, OUTPUT);
  _rows = rows;
  _columns = columns;
}

void LCD::begin() {
  _serial.begin(9600);                // 9600 baud is chip comm speed
  
  set_geometry(_rows, _columns);
  delay(500);                           // pause to allow LCD EEPROM to program

  set_brightness(0xFF);
  delay(1000);                          // pause to allow LCD EEPROM to program

  print("?s6");                  // set tabs to six spaces
  delay(1000);                          // pause to allow LCD EEPROM to program

  print("?D00000000000000000");  // define special characters
  delay(300);                           // delay to allow write to EEPROM
  // see moderndevice.com for a handy custom char generator (software app)
  clear_screen();
  delay(10);
  print("...");


  //crashes LCD without delay
  print("?D11010101010101010");
  delay(300);

  print("?D21818181818181818");
  delay(300);

  print("?D31c1c1c1c1c1c1c1c");
  delay(300);

  print("?D41e1e1e1e1e1e1e1e");
  delay(300);

  print("?D51f1f1f1f1f1f1f1f");
  delay(300);

  print("?D60000000000040E1F");
  delay(300);

  print("?D70000000103070F1F");
  delay(300);

  disable_cursor();
  delay(300);
}

void LCD::clear_screen() {
  print("?f");
}

void LCD::set_geometry(byte rows, byte columns) {
  print("?G%d%02d", rows, columns);
}

void LCD::set_brightness(byte brightness) {
  print("?B%02X", brightness);
}

void LCD::set_cursor(byte row, byte column) {
  set_cursor_row(row);
  set_cursor_column(column);
}

void LCD::set_cursor_row(byte row) {
  print("?y%d", row);
}

void LCD::set_cursor_column(byte column) {
  print("?x%02d", column);
}

void LCD::underline_cursor() {
  print("?c3");
}

void LCD::blink_cursor() {
  print("?c2");
}

void LCD::disable_cursor() {
  print("?c0");
}

void LCD::print(char const *fmt, ... ) {
  char buf[20];
  va_list args;
  va_start (args, fmt );
  vsnprintf(buf, 20, fmt, args);
  va_end (args);
  _serial.print(buf);
}
