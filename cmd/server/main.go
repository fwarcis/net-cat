package main

import (
	"bufio"
	"container/list"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"netcat/internal/common/comm"
	"netcat/internal/common/conf"
	"netcat/internal/common/errs"
	"netcat/internal/common/flags"
	"netcat/internal/server/history"
	"netcat/internal/server/sharing"
	"netcat/internal/server/users"
)

const historyPath = "assets/server/chat-history.txt"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})))

	host, port, cfgErr := AddressFromConfigAndArg()
	if cfgErr != nil {
		slog.Error(cfgErr.Error())
		return
	}

	listnr, listnrCreatErr := net.Listen("tcp", host+":"+port)
	if listnrCreatErr != nil {
		slog.Error(listnrCreatErr.Error())
		return
	}
	defer func() {
		err := listnr.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	fmt.Println("Listening on the port :" + port)

	welcomeMsg, welcErr := WelcomeMessage()
	if welcErr != nil {
		slog.Error(welcErr.Error())
		return
	}
	connCount := 0
	regedUsers := list.New()
	for {
		conn, acceptingErr := listnr.Accept()
		if acceptingErr != nil {
			slog.Error(acceptingErr.Error())
			continue
		}
		if connCount >= 10 {
			err := conn.Close()
			if err != nil {
				slog.Error(err.Error())
			}
			continue
		}
		connCount++

		go func() {
			defer func() {
				err := conn.Close()
				if err != nil {
					slog.Error(err.Error())
				}
			}()

			lnRecvr := comm.LineReceiver{Scanner: bufio.NewScanner(conn)}
			lnSender := comm.LineSender{bufio.NewWriter(conn)}
			user, userCreatErr := users.CreateUser(&lnRecvr, &lnSender, welcomeMsg)
			if userCreatErr != nil {
				_ = lnSender.Send("Error when creating a user.")
				slog.Error(userCreatErr.Error())
				return
			}
			userElem := regedUsers.PushBack(user)
			defer regedUsers.Remove(userElem)

			history, histErr := OpenHistory()
			if histErr != nil {
				slog.Error(histErr.Error())
			}
			defer sharing.NotifyLeft(regedUsers, userElem, history)

			histSendErr := sharing.SendHistory(history, user)
			if histSendErr != nil {
				slog.Error(histSendErr.Error())
				if errs.IsEOFLike(histSendErr) {
					return
				}
			}

			sharing.NotifyJoined(regedUsers, userElem, history)

			sendingErr := sharing.StreamMessages(
				regedUsers, userElem, &lnRecvr, history)
			if sendingErr != nil {
				slog.Error(sendingErr.Error())
				return
			}
		}()
	}
}

func WelcomeMessage() (string, error) {
	logo, err := os.ReadFile("assets/server/linux-logo.txt")
	if err != nil {
		return "", err
	}
	return "Welcome to TCP-Chat!\n" + string(logo), nil
}

func OpenHistory() (*history.History, error) {
	historyFile, err := os.OpenFile(
		historyPath,
		os.O_APPEND|os.O_RDWR|os.O_CREATE,
		0o700)
	if err != nil {
		return nil, err
	}
	return &history.History{
		Sender:   &comm.LineSender{bufio.NewWriter(historyFile)},
		Receiver: &comm.LineReceiver{Scanner: bufio.NewScanner(historyFile)},
	}, nil
}

func AddressFromConfigAndArg() (string, string, error) {
	cfg := conf.NewDefaultConfig()

	port := *flags.TextVar("port", &cfg.Port, "")
	flag.Parse()
	isPortSet := flags.IsSet("port")

	if flag.NArg() > 1 {
		return "", "", errors.New("[USAGE]: ./TCPChat $port")
	}
	if flag.NArg() == 1 && !isPortSet {
		err := port.UnmarshalText([]byte(flag.Arg(0)))
		if err != nil {
			return "", "", err
		}
	}
	return cfg.Host.String(), port.String(), nil
}
