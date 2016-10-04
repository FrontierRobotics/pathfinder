#include <SoftwareSerial.h>
#include <stdarg.h>

#define rxPin 30  // rxPin is immaterial - not used - just make this an unused Arduino pin number
#define txPin 10

SoftwareSerial serial =  SoftwareSerial(rxPin, txPin);


char N;
int I;
int ByteVar;

int NN;
int Remainder;
int Num_5;

void lcd_initialize() {
  pinMode(txPin, OUTPUT);
  serial.begin(9600);                // 9600 baud is chip comm speed

  lcd_set_geometry(4, 20);           // set display geometry,  4 x 20 characters in this case
  delay(500);                           // pause to allow LCD EEPROM to program

  lcd_set_brightness(0xFF);          // set backlight to ff hex, maximum brightness
  delay(1000);                          // pause to allow LCD EEPROM to program

  lcd_print("?s6");                  // set tabs to six spaces
  delay(1000);                          // pause to allow LCD EEPROM to program

  lcd_print("?D00000000000000000");  // define special characters
  delay(300);                           // delay to allow write to EEPROM
  // see moderndevice.com for a handy custom char generator (software app)
  lcd_clear();
  delay(10);
  lcd_print("...");


  //crashes LCD without delay
  lcd_print("?D11010101010101010");
  delay(300);

  lcd_print("?D21818181818181818");
  delay(300);

  lcd_print("?D31c1c1c1c1c1c1c1c");
  delay(300);

  lcd_print("?D41e1e1e1e1e1e1e1e");
  delay(300);

  lcd_print("?D51f1f1f1f1f1f1f1f");
  delay(300);

  lcd_print("?D60000000000040E1F");
  delay(300);

  lcd_print("?D70000000103070F1F");
  delay(300);

  lcd_disable_cursor();
  delay(300);
}

void lcd_clear() {
  lcd_print("?f");
}

void lcd_set_geometry(byte rows, byte columns) {
  lcd_print("?G%d%02d", rows, columns);
}

void lcd_set_brightness(byte brightness) {
  lcd_print("?B%02X", brightness);
}

void lcd_set_cursor(byte row, byte column) {
  lcd_set_cursor_row(row);
  lcd_set_cursor_column(column);
}

void lcd_set_cursor_row(byte row) {
  lcd_print("?y%d", row);
}

void lcd_set_cursor_column(byte column) {
  lcd_print("?x%02d", column);
}

void lcd_underline_cursor() {
  lcd_print("?c3");
}

void lcd_blink_cursor() {
  lcd_print("?c2");
}

void lcd_disable_cursor() {
  lcd_print("?c0");
}

void lcd_print(char const *fmt, ... ) {
  char buf[20];
  va_list args;
  va_start (args, fmt );
  vsnprintf(buf, 20, fmt, args);
  va_end (args);
  serial.print(buf);
}

