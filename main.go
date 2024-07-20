package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// processing state enums
const (
	WAIT_FOR_MSG = iota
	IN_MSG
)

func main(){
	port := flag.Int("port", 9090, "Specify the port for the TCP server. (default is 9090)")
	numOfWorkers := flag.Int("workers", 2, "Specify the number of workers in the thread pool. (default is 2)")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d",port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting the server: %v", err)
		os.Exit(1)
	}
}