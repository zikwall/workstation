package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	processId := getProcessIdent()

	toStdout := makeStdoutWriter(processId)
	sig := makeNotifier()
loop:
	for {
		select {
		case <-sig:
			break loop
		case <-time.After(time.Second * 2):
			toStdout()
		}
	}

	fmt.Println("Good Luck")
}

func getProcessIdent() string {
	var processId string
	flag.StringVar(&processId, "id", "undefined", "process identifier")
	flag.Parse()

	return processId
}

func makeNotifier() <-chan os.Signal {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	return sig
}

func makeStdoutWriter(processId string) func() {
	write := func(processId string) {
		_, _ = fmt.Fprintln(os.Stdout, fmt.Sprintf("go,process,19666,720,%s", processId)) // <- last is ID
	}

	return func() {
		write(processId)
	}
}
