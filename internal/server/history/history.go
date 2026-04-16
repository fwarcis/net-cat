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

func (h *History) Receive(
	beforeScan func(),
	afterScan func(ln string) bool,
) {
	h.Receiver.Receive(beforeScan, afterScan)
}

func (h *History) Err() error {
	return h.Receiver.Err()
}
