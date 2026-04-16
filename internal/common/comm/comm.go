package comm

import (
	"bufio"
	"io"
)

type Sender interface {
	Send(str string) error
}

type LineReceiver struct {
	Scanner   *bufio.Scanner
	hasNoMore bool
}

func (r *LineReceiver) Receive(
	beforeScan func(),
	afterScan func(ln string) bool,
) {
	if beforeScan == nil {
		for {
			r.hasNoMore = !r.Scanner.Scan()
			if r.hasNoMore || !afterScan(r.Scanner.Text()) {
				return
			}
		}
	} else {
		for {
			beforeScan()
			r.hasNoMore = !r.Scanner.Scan()
			if r.hasNoMore || !afterScan(r.Scanner.Text()) {
				return
			}
		}
	}
}

// bufio.Scanner.Err() or io.EOF
func (r *LineReceiver) Err() error {
	scanrErr := r.Scanner.Err()
	if r.hasNoMore && scanrErr == nil {
		return io.EOF
	}
	return scanrErr
}

type LineSender struct {
	Writer *bufio.Writer
}

func (s *LineSender) Send(str string) error {
	_, err := s.Writer.WriteString(str + "\n")
	if err != nil {
		return err
	}
	return s.Writer.Flush()
}
