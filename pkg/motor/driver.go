package motor

import "github.com/andycondon/pathfinder/pkg/status"

type Tx func(w []byte, r []byte) error

type ReadStatus func([]byte) (status.Reading, error)

type Driver struct {
	Tx
	ReadStatus
	Left, Right *Motor
}

func (d *Driver) Stop() (status.Reading, error) {
	return d.txPair(d.Left.Stop(), d.Right.Stop())
}

func (d *Driver) Forward(s Speed) (status.Reading, error) {
	return d.txPair(d.Left.Forward(s), d.Right.Forward(s))
}

func (d *Driver) Reverse(s Speed) (status.Reading, error) {
	return d.txPair(d.Left.Reverse(s), d.Right.Reverse(s))
}

func (d *Driver) RotateLeft(s Speed) (status.Reading, error) {
	return d.txPair(d.Left.Reverse(s), d.Right.Forward(s))
}

func (d *Driver) RotateRight(s Speed) (status.Reading, error) {
	return d.txPair(d.Left.Forward(s), d.Right.Reverse(s))
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
