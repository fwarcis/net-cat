package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})))

	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Println("nc $IP $port")
		return
	}

	conn, err := net.Dial("tcp", flag.Arg(0)+":"+flag.Arg(1))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	recv := bufio.NewScanner(conn)
	retrv := bufio.NewWriter(conn)

	isFinished := false
	wg := sync.WaitGroup{}
	wg.Go(func() { ReceiveAndPrintln(*recv, &isFinished) })
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
		slog.Error(err.Error())
	}
}

func InputAndRetrieve(w bufio.Writer, isFinished *bool) {
	for !*isFinished {
		input, rErr := bufio.NewReader(os.Stdin).ReadString('\n')
		if rErr != nil {
			slog.Error(rErr.Error())
		}

		_, wErr := w.Write([]byte(input))
		if wErr != nil {
			slog.Error(wErr.Error())
		}

		flshErr := w.Flush()
		if flshErr != nil {
			slog.Error(flshErr.Error())
		}
	}
}
