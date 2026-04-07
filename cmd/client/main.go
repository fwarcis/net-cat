package main

import (
	"bufio"
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

	conn, err := net.Dial("tcp", AddressFromConfigOrArgs())
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err.Error())
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
	for scanr.Scan() {
		text := scanr.Text()
		if text == "" {
			fmt.Print("\n")
		} else {
			fmt.Print(text)
		}
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

func AddressFromConfigOrArgs() string {
	cfg := conf.NewDefaultConfig()

	host := *flags.TextVar("host", &cfg.Host, "")
	port := *flags.TextVar("port", &cfg.Port, "")
	flag.Parse()
	isHostSet := flags.IsSet("host")
	isPortSet := flags.IsSet("port")

	if flag.NArg() > 2 {
		log.Fatalln("nc $IP $port")
	}
	if flag.NArg() >= 1 && !isHostSet {
		err := host.UnmarshalText([]byte(flag.Arg(0)))
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	if flag.NArg() == 2 && !isPortSet {
		err := port.UnmarshalText([]byte(flag.Arg(1)))
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	return host.String() + ":" + port.String()
}

func WaitForKeyboardInterrupt() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		os.Exit(130)
	}()
}
