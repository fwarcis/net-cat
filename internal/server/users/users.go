package users

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"netcat/internal/common/ansi"
	"netcat/internal/common/comm"
)

type TextGetter interface {
	Get(text string) error
}

type User struct {
	Sender comm.Sender
	Name   string

	lastTextLen int

	mu sync.Mutex
}

func (u *User) Get(text string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.Sender.Send(text)
}

func (u *User) prompt() string {
	now := time.Now().UTC().Format(time.DateTime)
	timeAndUsername := fmt.Sprintf("[%s][%s]:\n", now, u.Name)
	u.lastTextLen = len([]rune(timeAndUsername)) - 1
	return timeAndUsername + ansi.CarriageOffset(
		u.lastTextLen, -1)
}

func (u *User) BackCarriageAfterLF() error {
	return u.Get(ansi.CarriageOffset(u.lastTextLen, -2))
}

func (u *User) GetPrompt() error {
	return u.Get(u.prompt())
}

func (u *User) SendPromptedMessage(getter TextGetter, msg string) error {
	return getter.Get(u.prompt() + msg)
}

func (u *User) NotifyJoined(getter TextGetter) error {
	return getter.Get(u.Name + " has joined our chat...")
}

func (u *User) NotifyLeft(getter TextGetter) error {
	return getter.Get(u.Name + " has left our chat...")
}

func CreateUser(
	recvr *comm.LineReceiver,
	sender comm.Sender,
	welcomeMsg string,
) (*User, error) {
	const enterNameMessage = "[ENTER YOUR NAME]: "

	senderErr := sender.Send(welcomeMsg + enterNameMessage)
	if senderErr != nil {
		return nil, senderErr
	}

	var username string
	recvr.Receive(nil, func(input string) bool {
		if input != "" {
			username = input
			return false
		}
		senderErr = sender.Send(enterNameMessage)
		return senderErr == nil
	})

	err := errors.Join(recvr.Err(), senderErr)
	if err != nil {
		return nil, err
	}
	u := &User{Sender: sender, Name: username}
	_ = u.BackCarriageAfterLF()
	return u, nil
}
