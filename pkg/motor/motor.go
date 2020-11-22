package motor

const (
	forward = 0x01
	reverse = 0x00
)

type Motor struct {
	Addr      byte
	Slow, Med byte
}

func (m *Motor) speed(s Speed) byte {
	switch s {
	case Full:
		return 0xFF
	case Slow:
		return m.Slow
	case Medium:
		return m.Med
	default:
		return 0x00
	}
}

func (m *Motor) Stop() []byte {
	return []byte{m.Addr, forward, m.speed(Stop)}
}

func (m *Motor) Forward(s Speed) []byte {
	return []byte{m.Addr, forward, m.speed(s)}
}

func (m *Motor) Reverse(s Speed) []byte {
	return []byte{m.Addr, reverse, m.speed(s)}
}
