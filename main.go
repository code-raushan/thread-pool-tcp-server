package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

// processing state enums
const (
	WAIT_FOR_MSG = iota
	IN_MSG
)

func main() {
	port := flag.Int("port", 9090, "Specify the port for the TCP server. (default is 9090)")
	numOfWorkers := flag.Int("workers", 2, "Specify the number of workers in the thread pool. (default is 2)")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting the server: %v", err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Fprintf(os.Stdout, "TCP server running on port %d\n", *port)
	clientCh := make(chan net.Conn)

	var wg sync.WaitGroup

	for i := 0; i < *numOfWorkers; i++ {
		wg.Add(1)
		go worker(clientCh, &wg)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error accepting connection: %v", err)
			continue
		}

		clientCh <- client
		fmt.Printf("Client connected: %s\n", client.RemoteAddr().String())
	}
}

func worker(clientCh chan net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for client := range clientCh {
		serveConnection(client)
	}
}

func serveConnection(conn net.Conn) {
	defer conn.Close()

	// sending acknowledgment to the client
	if _, err := conn.Write([]byte("*")); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to client: %v\n", err)
		return
	}

	state := WAIT_FOR_MSG

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break // Client closed the connection
			}
			fmt.Fprintf(os.Stderr, "Error reading from client: %v\n", err)
			return
		}

		for i := 0; i < n; i++ {
			switch state {
			case WAIT_FOR_MSG:
				if buf[i] == '^' {
					state = IN_MSG
					fmt.Fprintf(os.Stdout, "In-Message State\n")
				}
			case IN_MSG:
				if buf[i] == '$' {
					state = WAIT_FOR_MSG
					fmt.Fprintf(os.Stdout, "Wait-For-Message State\n")

				} else {
					buf[i]++
					if _, err := conn.Write(buf[i : i+1]); err != nil {
						fmt.Fprintf(os.Stderr, "Error writing to client: %v\n", err)
						return
					}
				}
			}
		}
	}
}
