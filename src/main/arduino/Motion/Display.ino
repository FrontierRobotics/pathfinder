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

  lcd_print("?G420");                // set display geometry,  4 x 20 characters in this case
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

void lcd_set_brightness(byte brightness) {
  lcd_print("?B%02X", brightness);
}

void lcd_set_cursor(byte x, byte y) {
  lcd_print("?x%02d?y%d", x, y);
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

void lcd_demo() {
  serial.print("?f");                   // clear the LCD

  delay(1000);

  serial.print("?f");                   // clear the LCD
  delay(100);
  delay(3000);

  serial.print("?x00?y0");              // cursor to first character of line 0

  serial.print("LCD117 serial ?6?7");
  serial.print("?0?1?2?3?4?5");     // display special characters

  delay(3000);

  serial.print("?x00?y1");              // move cursor to beginning of line 1
  serial.print("moderndevice.com");     // crass commercial message

  delay(6000);                            // pause six secs to admire

  serial.print("?f");                   // clear the LCD

  serial.print("?x00?y0");              // move cursor to beginning of line 0

  serial.print(" LCD 117 chip by");     // displys LCD #117 on the screen


  serial.print("?x00?y1");              // cursor to first character of line 1
  serial.print(" phanderson.com");


  delay(3000);                            // pause three secs to admire

  serial.print("?f");                   // clear the screen

  serial.print("?x00?y0");              // locate cursor to beginning of line 0
  serial.print("DEC   HEX   ASCI");     // print labels
  delay(100);
  // simple printing demonstation
  for (N = 42; N <= 122; N++) {           // pick an arbitrary part of ASCII chart - change as you wish
    serial.print("?x00?y1");           // locate cursor to beginning of line 1



    serial.print(N, DEC);               // display N in decimal format
    serial.print("?t");                 // tab in

    serial.print(N, HEX);               // display N in hexidecimal format
    serial.print("?t");                 // tab in

    // glitches on ASCII 63 "?"
    if (N == '?') {
      serial.print("??");              // the "??" displays a single '?' - see Phanderson 117 docs
    }
    else {
      serial.write(N);           // display N as an ASCII character
    }

    serial.print("   ");                // display 3 spaces (blanks) as ASCII characters


    delay(500);
  }




  delay (1000);
  serial.print("?y0?x00");          // cursor to beginning of line 0
  delay(10);
  serial.print("?l");               // clear line; custom char. 1
  delay(10);
  serial.print(" Bar Graph Demo");
  delay(10);
  serial.print("?n");               // cursor to beginning of line 1 + clear line 1
  delay(500);

  // bar graph demo - increasing bar
  for ( N = 0; N <= 80; N++) {        // 16 chars * 5 bits each = 80
    serial.print("?y1?x00");       // cursor to beginning of line 1
    delay(10);

    Num_5 = N / 5;                   // calculate solid black tiles
    for (I = 1; I <= Num_5; I++) {
      serial.print("?5");         // print custom character 5 - solid block tiles
      delay(8);
    }

    Remainder = N % 5;               // % sign is modulo operator - calculates remainder
    // now print the remainder
    serial.print("?");       // first half of the custom character command; see end note
    serial.print(Remainder, DEC);  // prints the custom character equal to remainder
    delay(8);
  }

  delay(50);

  for ( N = 80; N >= 0; N--) {        // decreasing bar - 16 chars * 5 bits each
    serial.print("?y1?x00");       // cursor to beginning of line 1
    delay(14);

    Num_5 = N / 5;                   // calculate solid black tiles
    for (I = 1; I <= Num_5; I++) {
      serial.print("?5");         // print custom character 5 - solid block tiles
      delay(8);
    }

    Remainder = N % 5;               // % sign is modulo operator - calculates remainder
    // now print the remainder
    serial.print("?");             // first half of the custom character command
    serial.print(Remainder, DEC);  // prints the custom character equal to remainder
    delay(8);
  }

  delay(500);
  serial.print("?f");               // clears screen
  delay(50);
  serial.print(".");
  delay(60);
  serial.print("?y0?x00");          // cursor to beginning of line 0
  delay(250);
  serial.print(" .");
  delay(10);

  serial.print("?D0000000000000001F");   // define special characters
  delay(300);                              // delay to allow write to EEPROM
  //crashes LCD without delay
  serial.print("?y0?x00");               // cursor to beginning of line 0
  serial.print("  .");                   // dots for user feedback
  delay(10);

  serial.print("?D10000000000001F1F");
  delay(300);
  serial.print("?y0?x00");               // cursor to beginning of line 0
  serial.print(" . ");
  delay(10);

  serial.print("?D200000000001F1F1F");
  delay(300);
  serial.print("?y0?x00");               // cursor to beginning of line 0
  serial.print(".  ");
  delay(10);

  serial.print("?D3000000001F1F1F1F");
  delay(300);
  serial.print("?y0?x00");               // cursor to beginning of line 0
  serial.print(" .");
  delay(10);

  serial.print("?D40000001F1F1F1F1F");
  delay(300);
  serial.print("?y0?x00");               // cursor to beginning of line 0
  serial.print("  .");
  delay(10);

  serial.print("?D500001F1F1F1F1F1F");
  delay(300);
  serial.print("?y0?x00");               // cursor to beginning of line 0
  serial.print(" . ");
  delay(10);

  serial.print("?D6001F1F1F1F1F1F1F");
  delay(300);
  serial.print("?y0?x00");               // cursor to beginning of line 0
  serial.print(".  ");
  delay(10);

  serial.print("?D71F1F1F1F1F1F1F1F");
  delay(300);
  serial.print("?y0?x00");               // cursor to beginning of line 0
  serial.print(" . ");
  delay(10);

  serial.print("?c0");                   // turn cursor off
  delay(300
       );

  serial.print("?f");                    // clear the LCD

  delay(1000);

  serial.print("?y0?x00");               // cursor to beginning of line 0
  delay(10);
  serial.print("?l");                    // clear line
  delay(10);
  serial.print("  Vertical Bar ");
  delay(10);
  serial.print("?n");                    // cursor to beginning of line 1 + clear line 1
  serial.print("      Demo     ");
  delay(500);

  // vertical bar graph demo - increasing bar
  for ( N = 0; N <= 15; N++) {
    serial.print("?y1?x00");            // cursor to beginning of line 1
    delay(10);

    if (N < 8) {
      serial.print("?y0?x00 ");        // cursor to beginning of line 1 and writes blank space
      serial.print("?y1?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }
    else {
      serial.print("?y1?x00?7");       // cursor to beginning of line 1 and writes black character
      serial.print("?y0?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }

  }
  delay(100);

  for ( N = 15; N >= 0; N--) {
    serial.print("?y1?x00");            // cursor to beginning of line 1
    delay(10);

    if (N < 8) {
      serial.print("?y0?x00 ");        // cursor to beginning of line 1 and writes blank space
      serial.print("?y1?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }
    else {
      serial.print("?y1?x00?7");       // cursor to beginning of line 1 and writes black character
      serial.print("?y0?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }

  }
  delay(50);
  for ( N = 0; N <= 15; N++) {             // decreasing bar
    serial.print("?y1?x00");            // cursor to beginning of line 1
    delay(10);

    if (N < 8) {
      serial.print("?y0?x00 ");        // cursor to beginning of line 1 and writes blank space
      serial.print("?y1?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }
    else {
      serial.print("?y1?x00?7");       // cursor to beginning of line 1 and writes black character
      serial.print("?y0?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }

  }
  delay(100);

  for ( N = 15; N >= 0; N--) {
    serial.print("?y1?x00");            // cursor to beginning of line 1
    delay(10);

    if (N < 8) {
      serial.print("?y0?x00 ");        // cursor to beginning of line 1 and writes blank space
      serial.print("?y1?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }
    else {
      serial.print("?y1?x00?7");       // cursor to beginning of line 1 and writes black character
      serial.print("?y0?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }

  }
  delay(50);
  for ( N = 0; N <= 15; N++) {
    serial.print("?y1?x00");            // cursor to beginning of line 1
    delay(10);

    if (N < 8) {
      serial.print("?y0?x00 ");        // cursor to beginning of line 1 and writes blank space
      serial.print("?y1?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }
    else {
      serial.print("?y1?x00?7");       // cursor to beginning of line 1 and writes black character
      serial.print("?y0?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }

  }
  delay(100);

  for ( N = 15; N >= 0; N--) {             // 16 chars * 8 bits each = 80
    serial.print("?y1?x00");            // cursor to beginning of line 1
    delay(10);

    if (N < 8) {
      serial.print("?y0?x00 ");        // cursor to beginning of line 1 and writes blank space
      serial.print("?y1?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }
    else {
      serial.print("?y1?x00?7");       // cursor to beginning of line 1 and writes black character
      serial.print("?y0?x00");
      Remainder = (N % 8);               // % sign is modulo operator - calculates remainder
      // now print the remainder
      serial.print("?");               // first half of the custom character command
      serial.print(Remainder, DEC);    // prints the custom character equal to remainder
      delay(10);

    }

  }
  delay(1000);
}

