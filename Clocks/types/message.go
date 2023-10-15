package types

type Message struct {
	Msg         string
	From        int
	Clock       int
	VectorClock []int
}

func (m *Message) GetClock() int {
	return m.Clock
}

func (m *Message) SetClock(clock int) {
	m.Clock = clock
}

func (m *Message) GetVectorClock() []int {
	return m.VectorClock
}

func (m *Message) SetVectorClock(vector []int) {
	m.VectorClock = vector
}

func (m *Message) IncrementVectorClock(index int) {
	if index == -1 {
		index = len(m.VectorClock) - 1
	}
	m.VectorClock[index] += 1
}
