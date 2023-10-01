package main

import (
	"fmt"
	"net"
	"os"

	"github.com/abdelrhman-basyoni/godis/godis"
)

func server() {
	// listen to tcp
	l, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed while creating godis server- %v\n", err)

	}
	defer l.Close()

	//aof init
	aof := godis.AOF
	defer aof.Close()
	//read the database (aof)
	aof.Read()

	fmt.Println("server listening to port 6379")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed while accepting godis connections- %v\n", err)
			continue // Continue to accept more connections
		}

		// Use goroutine to handle each connection concurrently
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	writer := godis.NewBasicWriter(conn)
	fmt.Println("New connection established")

	for {
		resp, err := godis.NewRespReader(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: reader failed- %v\n", err)
			return // Terminate this goroutine
		}

		value, err := resp.Read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed while parsing resp- %v\n", err)
			return // Terminate this goroutine
		}

		res := godis.HandleValue(value)

		writer.Write(res)
	}
}
