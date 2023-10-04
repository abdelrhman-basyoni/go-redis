package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/abdelrhman-basyoni/godis/godis"
)

type AppConfig struct {
	appendOnly string
}

var appConfig = AppConfig{
	appendOnly: "no",
}

func flagsInit() {

	// Usage: -name John
	flag.StringVar(&appConfig.appendOnly, "appendonly", appConfig.appendOnly, "choose 'yes' or 'no'")

}
func main() {
	flagsInit()
	flag.Parse()

	fmt.Println(appConfig)
	server()

}

func server() {

	// set configuration
	if appConfig.appendOnly == "yes" {

		godis.Conf.SetAO(true)
	}

	if err := godis.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed while initializing godis- %v\n", err)
		os.Exit(1)
	}
	// listen to tcp
	l, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed while creating godis server- %v\n", err)

	}
	defer l.Close()

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
