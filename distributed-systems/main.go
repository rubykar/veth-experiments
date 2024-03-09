package main

import (
	"fmt"
	"github.com/rubykar/distributed-testing/node"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	node1 := &node.Node{Port: 9090, Stop: make(chan struct{})} // Assuming Stop channel is exported
	node2 := &node.Node{Port: 9091, Stop: make(chan struct{})} // Assuming Stop channel is exported

	go func() {
		fmt.Println("Node 1 Starting.......")
		node1.StartServer()
	}()

	go func() {
		fmt.Println("Node 2 Starting.......")
		node2.StartServer()
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	<-terminate // Block until termination signal is received
	node1.StopServer()
	node2.StopServer()
	fmt.Println("Servers stopped.")
}
