package motor

import "errors"

type Tx func(w []byte, r []byte) error

type Driver struct {
	Tx
	Left, Right *Motor
}

func (d *Driver) D(cmd Command) error {
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
		return errors.New("unknown motor movement")
	}
}

func (d *Driver) txPair(m1Cmd, m2Cmd []byte) error {
	err := d.tx(m1Cmd)
	if err != nil {
		return err
	}
	return d.tx(m2Cmd)
}

func (d *Driver) tx(cmd []byte) error {
	return d.Tx(cmd, nil)
}
