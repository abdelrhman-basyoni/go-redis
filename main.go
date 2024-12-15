package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/abdelrhman-basyoni/godis/godis"
	goresp "github.com/abdelrhman-basyoni/goresp"
)

type AppConfig struct {
	appendOnly string
}

var appConfig = AppConfig{
	appendOnly: "yes",
}

func flagsInit() {

	// Usage: -name John
	flag.StringVar(&appConfig.appendOnly, "appendonly", appConfig.appendOnly, "choose 'yes' or 'no'")

}
func main() {
	defer recoverPanic()
	flagsInit()
	flag.Parse()

	fmt.Println(appConfig)
	server()

}

func server() {
	defer recoverPanic()
	// set configuration
	if appConfig.appendOnly == "yes" {

		godis.Conf.SetAO(true)
	}
	cl, err := godis.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed while initializing godis- %v\n", err)
		os.Exit(1)
	}

	defer cl()
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
	defer recoverPanic()
	defer conn.Close()
	writer := goresp.NewBasicWriter(conn)
	fmt.Println("New connection established")

	for {
		resp := goresp.NewRespReader(conn)

		value, err := resp.Read()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed while parsing resp- %v\n", err)
			return // Terminate this goroutine
		}

		res := godis.HandleValue(value)

		writer.Write(res)
	}
}

func recoverPanic() {
	if err := recover(); err != nil {
		fmt.Printf("RECOVERED: %v\n", err)
	}
}
