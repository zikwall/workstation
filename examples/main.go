package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/zikwall/workstation"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	ctx, cancelFunc := context.WithCancel(context.Background())

	station := workstation.BuildWorkstation(ctx, &GoSubprocessMonitor{})

	sampleFunctionC := func(a int, b string, t time.Time, o string) {
		fmt.Println(t.String(), a, b, o)
	}

	err := station.PerformAsync("first_process", workstation.Payload{"a": 1, "b": "6", "c": sampleFunctionC})

	if err != nil {
		log.Fatal(err)
	}

	err = station.PerformAsync("second_process", workstation.Payload{"a": 3, "b": "3", "c": sampleFunctionC})

	if err != nil {
		log.Fatal(err)
	}

	<-sig

	cancelFunc()

	if err := station.Shutdown(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Terminated all")
}

type GoSubprocessMonitor struct{}

func (m *GoSubprocessMonitor) Perform(instance workstation.Instantiable, key string, payload workstation.Payload) {
	// Note: check is successfully cast!
	a := payload["a"].(int)
	b := payload["b"].(string)
	c := payload["c"].(func(a int, b string, t time.Time, o string))

	watcher, err := executeGoSubprocess(key)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(instance.ProvideExecutionContext())

	defer func() {
		if err := watcher.w.Close(); err != nil {
			log.Fatal(err)
		}

		if err := watcher.cmd.Process.Kill(); err != nil && !isAlreadyFinished(err) {
			errorFailedkillprocess(key, watcher.cmd.Process.Pid, err)
		} else {
			infoSuccessKillprocess(key, watcher.cmd.Process.Pid)
		}

		cancel()
	}()

	goSubprocessOutput := make(chan string, 100)
	goSubprocessKilled := make(chan error, 1)

	// Runs a separate sub-thread, because when running in a single thread,
	// there is a lock while waiting for the buffer to be read.
	// In turn blocking by the reader will not allow the background task to finish gracefully
	go func() {
		bufioReader := bufio.NewReader(watcher.r)

		for {
			line, isPrefix, err := bufioReader.ReadLine()

			if err != nil {
				errorCloseProcessStdoutReader(key, err)

				return
			}

			str := string(line)

			if isPrefix || str == "" {
				continue
			}

			goSubprocessOutput <- str
		}
	}()

	// We listen to the Go subprocess termination signal,
	// this will provide an opportunity to remove the task from the pool and restart it if necessary
	//
	// Note: We listen to the context so as not to leave active goroutines when the task is completed
	go func() {
		select {
		case goSubprocessKilled <- watcher.cmd.Wait():
			return
		case <-ctx.Done():
			return
		}
	}()

	isCanceled := instance.GetIsCancelledChannel(key)

	for instance.ObserveProcessAlive(key) {
		select {
		case <-isCanceled:
			return
		case <-ctx.Done():
			return
		case err := <-goSubprocessKilled:
			errorProcessKilled(key, watcher.cmd.Process.Pid, err)
			return
		case outPartials := <-goSubprocessOutput:
			c(a, b, time.Now(), outPartials)
		}
	}
}
