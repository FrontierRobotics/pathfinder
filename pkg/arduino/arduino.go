package arduino

import (
	"github.com/andycondon/pathfinder/pkg/ir"
)

const (
	I2CAddress    = 0x1A
	StatusAddress = 0x10
)

type Status struct {
	IR ir.Reading
}

type Arduino struct {
	Tx
	IRArray *ir.SensorArray
}

type Tx func(w []byte, r []byte) error

func New(tx Tx, irArray *ir.SensorArray) *Arduino {
	return &Arduino{
		Tx:      tx,
		IRArray: irArray,
	}
}

func (a *Arduino) GetStatus() (Status, error) {
	read := make([]byte, 3)
	if err := a.Tx([]byte{StatusAddress}, read); err != nil {
		return Status{}, err
	}
	irSensor, err := a.IRArray.Reading(read)
	if err != nil {
		// Not possible to hit this error if read slice is correct size above.
		return Status{}, err
	}
	return Status{
		IR: irSensor,
	}, nil
}
