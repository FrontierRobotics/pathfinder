package motor

import (
	"errors"

	"github.com/andycondon/pathfinder/pkg/status"
)

type Tx func(w []byte, r []byte) error

type ReadStatus func([]byte) (status.Reading, error)

type Driver struct {
	Tx
	ReadStatus
	Left, Right *Motor
}

func (d *Driver) D(cmd Command) (status.Reading, error) {
	s := cmd.S
	switch cmd.M {
	case Park:
		return d.txPair(d.Left.Stop(), d.Right.Stop())
	case Forward:
		return d.txPair(d.Left.Forward(s), d.Right.Forward(s))
	case Reverse:
		return d.txPair(d.Left.Reverse(s), d.Right.Reverse(s))
	case RotateLeft:
		return d.txPair(d.Left.Reverse(s), d.Right.Forward(s))
	case RotateRight:
		return d.txPair(d.Left.Forward(s), d.Right.Reverse(s))
	default:
		return status.Reading{}, errors.New("unknown motor movement")
	}
}

func (d *Driver) txPair(m1Cmd, m2Cmd []byte) (status.Reading, error) {
	s, err := d.tx(m1Cmd)
	if err != nil {
		return s, err
	}
	return d.tx(m2Cmd)
}

func (d *Driver) tx(cmd []byte) (status.Reading, error) {
	read := make([]byte, 3)
	if err := d.Tx(cmd, read); err != nil {
		return status.Reading{}, err
	}
	return d.ReadStatus(read)
}
