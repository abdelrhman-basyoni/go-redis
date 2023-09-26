package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/abdelrhman-basyoni/godis/godis"
	resp_basic "github.com/abdelrhman-basyoni/godis/resp"
)

func server() {
	// listen to tcp
	l, err := net.Listen("tcp", ":6379")
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
	writer := resp_basic.NewBasicWriter(conn)
	for {

		resp := resp_basic.NewBasicReader(conn)

		value, err := resp.Read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed while parsing resp- %v\n", err)

		}
		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := make([]Value, len(value.Array[1:]))
		// Convert resp_basic.Value to Value
		for i, v := range value.Array[1:] {
			args[i] = convertToValue(v)

		}

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(resp_basic.Value{Typ: "string", Str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(convertToGodisValue(value))
		}

		result := handler(args)

		writer.Write(convertToResp(result))
	}
}

func convertToValue(rbValue resp_basic.Value) Value {
	return Value{
		typ:  rbValue.Typ,
		str:  rbValue.Str,
		num:  rbValue.Num,
		bulk: rbValue.Bulk,
		// Handle the array field if needed
	}
}

func convertToResp(rbValue Value) resp_basic.Value {
	return resp_basic.Value{
		Typ:  rbValue.typ,
		Str:  rbValue.str,
		Num:  rbValue.num,
		Bulk: rbValue.bulk,
		// Handle the array field if needed
	}
}

func convertToGodisValue(rbValue resp_basic.Value) godis.Value {
	return godis.Value{
		Typ:  rbValue.Typ,
		Str:  rbValue.Str,
		Num:  rbValue.Num,
		Bulk: rbValue.Bulk,
		// Handle the array field if needed
	}
}
