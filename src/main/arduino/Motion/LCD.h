#ifndef LCD_h
#define LCD_h

#include <Arduino.h>
#include <SoftwareSerial.h>

class LCD {
  public:
    LCD(int txPin, byte rows, byte columns);
    void begin();
    void clear_screen();
    void set_brightness(byte brightness);
    void set_cursor(byte row, byte column);
    void set_cursor_row(byte row);
    void set_cursor_column(byte column);
    void disable_cursor();
    void print(char const *fmt, ... );
  private:
    SoftwareSerial _serial;
    byte _rows, _columns;
    void set_geometry(byte rows, byte columns);
};

#endif
