package comm

import (
	"bufio"
	"errors"
)

type Sender interface {
	Send(str string) error
}

type LineReceiver struct {
	Scanner *bufio.Scanner
}

func (r *LineReceiver) Receive(before func(), yield func(ln string) bool) error {
	if before == nil {
		before = func() {}
	}

	for {
		before()

		if !r.Scanner.Scan() || !yield(r.Scanner.Text()) {
			break
		}
	}
	return r.Scanner.Err()
}

type LineSender struct {
	Writer *bufio.Writer
}

func (s *LineSender) Send(str string) error {
	_, err := s.Writer.WriteString(str + "\n")
	if err != nil {
		return err
	}

	err = s.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

type LineSenderWithHistory struct {
	Sender Sender

	HistorySender   *LineSender
	HistoryReceiver *LineReceiver
}

func (s *LineSenderWithHistory) Send(str string) error {
	err := s.HistorySender.Send(str)
	return errors.Join(err, s.Sender.Send(str))
}

func (s *LineSenderWithHistory) ReadAllHistory(yield func(string) bool) error {
	return s.HistoryReceiver.Receive(nil, yield)
}
