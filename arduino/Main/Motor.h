#ifndef Motor_h
#define Motor_h

#include <Arduino.h>

class Motor
{
public:
    Motor(int dPin, int pwmPin);
    void begin();
    void forward(byte speed);
    void reverse(byte speed);
    bool isReversed();
    bool isStopped();
    byte speed();
    const char* status();

private:
    byte _dPin, _pwmPin, _speed;
    bool _reversed;
};

#endif