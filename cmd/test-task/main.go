package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/ruslan-onishchenko/go-test-task/internal"
)

func main() {
	defaultTimeout := flag.Int("defaultTimeout", 3, "during this timeout, the request will try to get a message from the queue")
	port := flag.Int("port", 8080, "port where server will be started")
	maxQueuesCount := flag.Int("maxQueuesCount", 3, "max queues count")
	maxMessagesInOneQueueCount := flag.Int("maxMessagesInOneQueueCount", 3, "max messages in one queue count")

	flag.Parse()

	mux := http.NewServeMux()

	userH := internal.NewHandler(*maxQueuesCount, *maxMessagesInOneQueueCount, *defaultTimeout)
	mux.Handle("/queue/", userH)

	for range make([]struct{}, internal.ReconnectAttemptsCount) {
		addr := fmt.Sprintf("localhost:%d", *port)
		fmt.Println("server starting at ", addr)
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			println(err.Error())
		}

		fmt.Printf("recconect after %d seconds \n", internal.ReconnectTimerSec)
		time.Sleep(internal.ReconnectTimerSec * time.Second)
	}

	fmt.Printf("ReconnectAttemptsCount exided, count attempts was %d \n", internal.ReconnectAttemptsCount)
}
