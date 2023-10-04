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

func Init() (func(), error) {

	if Conf.ao {
		//aof init
		aof := AOF
		Conf.SetAO(false)
		defer Conf.SetAO(true)
		//read the database (aof)
		if err := aof.Read(); err != nil {
			return nil, err
		}
	}

	return func() {
		AOF.Close()
	}, nil
}
