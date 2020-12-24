package ir

type Tx func(w []byte, r []byte) error

type Reader struct {
	Tx
	Addr    byte
	IRArray *SensorArray
}

func (r *Reader) Get() (Reading, error) {
	read := make([]byte, 3)
	if err := r.Tx([]byte{r.Addr}, read); err != nil {
		return Reading{}, err
	}
	return r.ReadStatus(read)
}

func (r *Reader) ReadStatus(read []byte) (Reading, error) {
	irSensor, err := r.IRArray.Reading(read)
	if err != nil {
		// Not possible to hit this error if read slice is correct size above.
		return Reading{}, err
	}
	return irSensor, nil
}
