package main

import (
	"bufio"
	"container/list"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"netcat/internal/common/comm"
	"netcat/internal/common/conf"
	"netcat/internal/common/flags"
	"netcat/internal/server/history"
	"netcat/internal/server/sharing"
	"netcat/internal/server/users"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	host, port, err := AddressFromConfigAndArg()
	if err != nil {
		log.Println(err.Error())
		return
	}

	listnr, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		err := listnr.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	fmt.Println("Listening on the port :" + port)

	welcomeMsg, err := WelcomeMessage()
	if err != nil {
		log.Println(err.Error())
		return
	}
	connCount := 0
	regedUsers := list.New()
	for {
		conn, err := listnr.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if connCount >= 10 {
			err = conn.Close()
			if err != nil {
				log.Println(err.Error())
			}
			continue
		}
		connCount++

		go func() {
			defer func() {
				err := conn.Close()
				if err != nil {
					log.Println(err.Error())
				}
			}()

			lnRecvr := comm.LineReceiver{bufio.NewScanner(conn)}
			lnSender := comm.LineSender{bufio.NewWriter(conn)}
			user, err := users.CreateUser(&lnRecvr, &lnSender, welcomeMsg)
			if err != nil {
				log.Println(err.Error())
				return
			}
			userElem := regedUsers.PushBack(user)
			defer regedUsers.Remove(userElem)

			history, err := OpenHistory()
			if err != nil {
				log.Println(err.Error())
			}
			defer sharing.NotifyLeft(regedUsers, userElem, history)
			err = history.Receive(func(text string) bool {
				if text == "" {
					return true
				}

				err = user.Get(text + "\n")
				if err != nil {
					log.Println(err.Error())
				}
				return true
			})
			if err != nil {
				log.Println(err.Error())
			}

			sharing.NotifyJoined(regedUsers, userElem, history)

			err = sharing.SendMessages(regedUsers, userElem, &lnRecvr, history)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	}
}

func WelcomeMessage() (string, error) {
	logo, err := os.ReadFile("assets/server/linux-logo.txt")
	if err != nil {
		return "", err
	}
	return "Welcome to TCP-Chat!\n\n" + string(logo), nil
}

func OpenHistory() (*history.History, error) {
	historyFile, err := os.OpenFile(
		"assets/server/chat-history.txt",
		os.O_APPEND|os.O_RDWR|os.O_CREATE,
		0o700,
	)
	if err != nil {
		return nil, err
	}
	return &history.History{
		Sender:   &comm.LineSender{bufio.NewWriter(historyFile)},
		Receiver: &comm.LineReceiver{bufio.NewScanner(historyFile)},
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
