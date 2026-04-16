package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"netcat/internal/common/conf"
	"netcat/internal/common/flags"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	go WaitForKeyboardInterrupt()

	// addr, err := AddressFromConfigOrArgs()
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	conn, err := net.Dial("tcp", os.Args[1])
	// conn, err := net.Dial("tcp, addr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	recv := bufio.NewScanner(conn)
	retrv := bufio.NewWriter(conn)

	isFinished := false
	wg := sync.WaitGroup{}

	wg.Go(func() { ReceiveAndPrintln(*recv, &isFinished) })
	// without wg; to exit without waiting for input
	go InputAndRetrieve(*retrv, &isFinished)

	wg.Wait()
}

func ReceiveAndPrintln(scanr bufio.Scanner, isFinished *bool) {
	if scanr.Scan() {
		fmt.Print(scanr.Text())
	}
	for scanr.Scan() {
		fmt.Print("\n" + scanr.Text())
	}

	*isFinished = true
	err := scanr.Err()
	if err != nil {
		log.Println(err.Error())
	}
}

func InputAndRetrieve(w bufio.Writer, isFinished *bool) {
	for !*isFinished {
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println(err.Error())
		}

		_, err = w.Write([]byte(input))
		if err != nil {
			log.Println(err.Error())
		}

		err = w.Flush()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func AddressFromConfigOrArgs() (string, error) {
	cfg := conf.NewDefaultConfig()

	host := *flags.TextVar("host", &cfg.Host, "")
	port := *flags.TextVar("port", &cfg.Port, "")
	flag.Parse()
	isHostSet := flags.IsSet("host")
	isPortSet := flags.IsSet("port")

	if flag.NArg() > 2 {
		return "", errors.New("nc $IP $port")
	}
	if flag.NArg() >= 1 && !isHostSet {
		err := host.UnmarshalText([]byte(flag.Arg(0)))
		if err != nil {
			return "", err
		}
	}
	if flag.NArg() == 2 && !isPortSet {
		err := port.UnmarshalText([]byte(flag.Arg(1)))
		if err != nil {
			return "", err
		}
	}
	return host.String() + ":" + port.String(), nil
}

func WaitForKeyboardInterrupt() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		os.Exit(0)
	}()
}
