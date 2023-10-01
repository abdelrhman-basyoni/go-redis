package godis

import (
	"fmt"
	"strings"
)

func HandleValue(value Value) []byte {
	errRes := Value{typ: "string", str: ""}.Marshal()

	if value.typ != "array" {
		return errRes

	}

	if len(value.array) == 0 {
		return errRes

	}

	command := strings.ToUpper(value.array[0].bulk)
	handleAOF(command, value)
	args := value.array[1:]
	handler, ok := Handlers[command]

	if !ok {
		fmt.Println("Invalid command: ", command)
		return errRes

	}

	result := handler(args)

	var bytes = result.Marshal()

	return bytes
}

func handleAOF(command string, value Value) {
	if command == "SET" || command == "HSET" {
		AOF.Write(value)
	}
}
