package sharing

import (
	"container/list"

	"netcat/internal/common/comm"
	"netcat/internal/common/errs"
	"netcat/internal/server/history"
	"netcat/internal/server/users"
)

func SendHistory(hist *history.History, user *users.User) error {
	var userMsgErr error
	hist.Receive(nil, func(ln string) bool {
		if ln == "" {
			return true
		}

		userMsgErr = user.Get(ln)
		return !errs.IsEOFLike(userMsgErr)
	})
	_ = hist.Err()
	return userMsgErr
}

func NotifyLeft(
	regedUsers *list.List,
	userElem *list.Element,
	hist *history.History,
) {
	senderUser := userElem.Value.(*users.User)
	for u := regedUsers.Front(); u != nil; u = u.Next() {
		if u == userElem {
			continue
		}

		recvrUser := u.Value.(*users.User)
		_ = senderUser.NotifyLeft(recvrUser)
		_ = recvrUser.GetPrompt()
	}
	_ = senderUser.NotifyLeft(hist)
	_ = hist.Err()
}

func NotifyJoined(
	regedUsers *list.List,
	userElem *list.Element,
	hist *history.History,
) {
	senderUser := userElem.Value.(*users.User)
	for u := regedUsers.Front(); u != nil; u = u.Next() {
		if u == userElem {
			continue
		}

		recvrUser := u.Value.(*users.User)
		_ = senderUser.NotifyJoined(recvrUser)
		_ = recvrUser.GetPrompt()
	}
	_ = senderUser.NotifyJoined(hist)
	_ = hist.Err()
}

func StreamMessages(
	regedUsers *list.List,
	senderUserElem *list.Element,
	recvr *comm.LineReceiver,
	hist *history.History,
) error {
	senderUser := senderUserElem.Value.(*users.User)

	// _ = senderUser.BackCarriageAfterLF()
	before := func() {
		_ = senderUser.GetPrompt()
	}
	recvr.Receive(before, func(ln string) bool {
		_ = senderUser.BackCarriageAfterLF()
		if ln == "" {
			return true
		}

		for uElem := regedUsers.Front(); uElem != nil; uElem = uElem.Next() {
			if uElem == senderUserElem {
				continue
			}

			go func() {
				recvrUser := uElem.Value.(*users.User)
				_ = senderUser.SendPromptedMessage(recvrUser, ln)
				_ = recvrUser.GetPrompt()
			}()
		}
		_ = senderUser.SendPromptedMessage(hist, ln)
		_ = hist.Err()

		return true
	})
	return recvr.Err()
}
