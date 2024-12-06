package godis

import (
	"fmt"
	"strings"
)

func HandleValue(value Value) []byte {
	errRes := Value{Typ: "string", Str: ""}.Marshal()

	if value.Typ != "Array" {
		return errRes

	}

	if len(value.Array) == 0 {
		return errRes

	}

	command := strings.ToUpper(value.Array[0].Bulk)

	args := value.Array[1:]
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
