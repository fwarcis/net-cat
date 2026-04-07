package sharing

import (
	"container/list"

	"netcat/internal/common/comm"
	"netcat/internal/server/history"
	"netcat/internal/server/users"
)

func NotifyLeft(
	regedUsers *list.List,
	userElem *list.Element,
	history *history.History,
) {
	senderUser := userElem.Value.(*users.User)
	for u := regedUsers.Front(); u != nil; u = u.Next() {
		if u == userElem {
			continue
		}

		recvrUser := u.Value.(*users.User)
		_ = recvrUser.Get("")
		_ = senderUser.NotifyLeft(recvrUser)
		_ = recvrUser.Get("")
		_ = recvrUser.SendMessage(recvrUser, "")
	}
	_ = history.Get("")
	_ = senderUser.NotifyLeft(history)
}

func NotifyJoined(
	regedUsers *list.List,
	userElem *list.Element,
	history *history.History,
) {
	senderUser := userElem.Value.(*users.User)
	for u := regedUsers.Front(); u != nil; u = u.Next() {
		if u == userElem {
			continue
		}

		recvrUser := u.Value.(*users.User)
		_ = recvrUser.Get("")
		_ = senderUser.NotifyJoined(recvrUser)
		_ = recvrUser.Get("")
		_ = recvrUser.SendMessage(recvrUser, "")
	}
	_ = history.Get("")
	_ = senderUser.NotifyJoined(history)
}

func SendMessages(
	regedUsers *list.List,
	userElem *list.Element,
	recvr *comm.LineReceiver,
	history *history.History,
) error {
	senderUser := userElem.Value.(*users.User)

	before := func() {
		_ = senderUser.SendMessage(senderUser, "")
	}
	return recvr.Receive(before, func(text string) bool {
		if text == "" {
			return true
		}

		for u := regedUsers.Front(); u != nil; u = u.Next() {
			if u == userElem {
				continue
			}

			recvrUser := u.Value.(*users.User)
			_ = recvrUser.Get("")
			_ = senderUser.SendMessage(recvrUser, text+"\n")
			_ = recvrUser.SendMessage(recvrUser, "")
		}
		_ = history.Get("")
		_ = senderUser.SendMessage(history, text)

		return true
	})
}
