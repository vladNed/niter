//go:build js && wasm
// +build js,wasm

package main

import (
	"log"
	"time"

	"syscall/js"

	"github.com/indexone/niter/core/discovery"
)

const VERSION = "0.0.1"

func main() {
	js.Global().Set("wasmVersion", VERSION)

	wsConn := initialize()
	defer wsConn.Close()

	wsConn.Start()

	// This is a blocking call to keep the wasm running.
	<-make(chan struct{})
}

// This method initialized all resources needed for the P2P node and server
// communication.
func initialize() *discovery.WSClient {
	log.Println("Initializing P2P node.")
	for {
		wsConn, err := discovery.NewWSClient()
		if err != nil {
			retryDelay()
			continue
		}

		return wsConn
	}
}

// This method is used to retry the initialization of the node in case of failure.
func retryDelay() {
	retryDelay := 2 * time.Second
	log.Println("Retrying initializing node")
	time.Sleep(retryDelay)
}
