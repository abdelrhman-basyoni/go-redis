package godis

import (
	"fmt"
	"strings"
)

func HandleValue(value Value) []byte {
	errRes := Value{Typ: "string", Str: "Invalid command"}

	if value.Typ != "array" {
		return errRes.Marshal()

	}

	if len(value.Array) == 0 {
		return errRes.Marshal()

	}

	command := strings.ToUpper(value.Array[0].Bulk)

	args := value.Array[1:]
	handler, ok := Handlers[command]

	if !ok {
		fmt.Println("Invalid command: ", command)
		errRes.Str = fmt.Sprintf("Invalid command: %v", command)
		return errRes.Marshal()

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
