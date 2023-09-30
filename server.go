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
	fmt.Println("server listening to port 6379")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed while creating godis server- %v\n", err)

	}
	defer l.Close()

	//aof init
	aof, err := godis.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()
	//read the database (aof)
	// aof.Read(func(value Value) bool {
	// 	command := strings.ToUpper(value.array[0].bulk)
	// 	args := value.array[1:]

	// 	handler, ok := Handlers[command]
	// 	if !ok {
	// 		fmt.Println("Invalid command: ", command)
	// 		return false
	// 	}

	// 	handler(args)
	// 	return true
	// })
	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed while accepting godis connections- %v\n", err)
	}
	defer conn.Close()
	writer := godis.NewBasicWriter(conn)
	for {

		resp, err := godis.NewRespReader(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: reader failed- %v\n", err)

		}

		value, err := resp.Read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed while parsing resp- %v\n", err)

		}

		res := godis.HandleValue(value)

		writer.Write(res)
	}
}
