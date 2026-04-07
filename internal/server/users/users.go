package users

import (
	"errors"
	"fmt"
	"time"

	"netcat/internal/common/comm"
)

type TextGetter interface {
	Get(text string) error
}

type User struct {
	Sender comm.Sender
	Name   string
}

func (u *User) SendMessage(getter TextGetter, msg string) error {
	now := time.Now().UTC().Format(time.DateTime)
	return getter.Get(fmt.Sprintf("[%s][%s]:%s", now, u.Name, msg))
}

func (u *User) NotifyJoined(getter TextGetter) error {
	return getter.Get(u.Name + " has joined our chat...")
}

func (u *User) NotifyLeft(getter TextGetter) error {
	return getter.Get(u.Name + " has left our chat...")
}

func (u *User) Get(text string) error {
	return u.Sender.Send(text)
}

func CreateUser(
	recvr *comm.LineReceiver,
	sender comm.Sender,
	welcomeMsg string,
) (_ *User, err error) {
	const enterNameMessage = "[ENTER YOUR NAME]: "

	err = sender.Send(welcomeMsg + enterNameMessage)
	if err != nil {
		return nil, err
	}

	var username string
	recvrErr := recvr.Receive(nil, func(input string) bool {
		if input != "" {
			username = input
			return false
		}
		senderErr := sender.Send(enterNameMessage)
		if senderErr != nil {
			err = senderErr
			return false
		}
		return true
	})
	err = errors.Join(err, recvrErr)
	if err != nil {
		return nil, err
	}

	return &User{sender, username}, nil
}
