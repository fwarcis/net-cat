package history

import (
	"netcat/internal/common/comm"
)

type History struct {
	Sender   comm.Sender
	Receiver *comm.LineReceiver
}

func (h *History) Get(text string) error {
	return h.Sender.Send(text)
}

func (h *History) Receive(yield func(string) bool) error {
	return h.Receiver.Receive(nil, yield)
}
