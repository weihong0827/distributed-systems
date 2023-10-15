package types

type Clock interface {
	GetClock() int
	SetClock(int)
	GetVectorClock() []int
	SetVectorClock([]int)
	IncrementVectorClock(int)
}

func UpdateClockWithMsg(entity Clock, msg Message, index int) {
	if index == -1 {
		index = len(entity.GetVectorClock()) - 1
	}
	if msg.Clock > entity.GetClock() {
		entity.SetClock(msg.Clock + 1)
	} else {
		entity.SetClock(entity.GetClock() + 1)
	}
	UpdateVectorClockWithMsg(entity, msg, index)
}

func UpdateVectorClockWithMsg(entity Clock, msg Message, index int) {
	localVectorClock := entity.GetVectorClock()
	messageVectorClock := msg.GetVectorClock()

	for i := range localVectorClock {
		// Take the maximum value for each entry in the vector clocks
		if messageVectorClock[i] > localVectorClock[i] {
			localVectorClock[i] = messageVectorClock[i]
		}
	}

	entity.IncrementVectorClock(index)

	entity.SetVectorClock(localVectorClock)
}

func UpdateMsgWithClock(entity Clock, message Message) {
	message.SetClock(entity.GetClock())
	message.SetVectorClock(entity.GetVectorClock())
}

func CheckCausalViolation(entity Clock, message Message) bool {
	// Check for at least one grater than
	atLeastOneGreater := false
	for i, val := range entity.GetVectorClock() {
		if val < message.GetVectorClock()[i] {
			atLeastOneGreater = true
		}
		if val > message.GetVectorClock()[i] {
			return false
		}
	}
	return atLeastOneGreater
}
