#include <Arduino.h>
#include "Motor.h"

Motor::Motor(int dPin, int pwmPin) : _dPin(dPin), _pwmPin(pwmPin), _speed(0), _reversed(false)
{
    pinMode(dPin, OUTPUT);
}

void Motor::begin()
{
    forward(0x00);
}

void Motor::forward(byte speed)
{
    _reversed = false;
    _speed = speed;
    digitalWrite(_dPin, HIGH);
    analogWrite(_pwmPin, speed);
}

void Motor::reverse(byte speed)
{
    _reversed = true;
    _speed = speed;
    digitalWrite(_dPin, LOW);
    analogWrite(_pwmPin, speed);
}

bool Motor::isReversed()
{
    return _reversed;
}

bool Motor::isStopped()
{
    return 0x00 == _speed;
}

byte Motor::speed()
{
    return _speed;
}

const char *Motor::status()
{
    if (isStopped())
    {
        return "STOP";
    }
    char *result = (char *)malloc(sizeof(char) * 4);
    if (isReversed())
    {
        sprintf(result, "R%03d", speed());
    }
    else
    {
        sprintf(result, "F%03d", speed());
    }
    return result;
}