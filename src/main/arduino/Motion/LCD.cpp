#include <Arduino.h>
#include <SoftwareSerial.h>
#include "LCD.h"

#define LCD_RX_PIN 30  // rxPin is immaterial - not used - just make this an unused Arduino pin number
#define BAUD_RATE 9600 // 9600 baud is chip comm speed

LCD::LCD(int txPin, byte rows, byte columns) : _serial(SoftwareSerial(LCD_RX_PIN, txPin)) {
  pinMode(txPin, OUTPUT);
  _rows = rows;
  _columns = columns;
}

void LCD::begin() {
  _serial.begin(BAUD_RATE);

  set_geometry(_rows, _columns);
  set_brightness(0xFF);
  set_tabs(6);

  define_character("?D00000000000000000");
  clear_screen();
  delay(10);
  print("...");

  disable_cursor();
  delay(300);
}

void LCD::clear_screen() {
  print("?f");
}

void LCD::set_brightness(byte brightness) {
  print("?B%02X", brightness);
  eeprom_delay(1000);
}

void LCD::set_tabs(byte tabs) {
  print("?s%d", tabs);
  eeprom_delay(1000);
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

//Private Methods

void LCD::set_geometry(byte rows, byte columns) {
  print("?G%d%02d", rows, columns);
  eeprom_delay(500);
}

// see moderndevice.com for a handy custom char generator (software app)
void LCD::define_character(const char* definition) {
  print(definition);
  eeprom_delay(300);
}

void LCD::eeprom_delay(unsigned long milliseconds) {
  delay(milliseconds);
}
