package godis

import (
	"fmt"
	"strconv"
	"time"
)

type commandFunction func([]Value) Value

var Handlers = map[string]commandFunction{
	"PING":   ping,
	"SET":    set,
	"GET":    get,
	"HSET":   hset,
	"HGET":   hget,
	"DEL":    del,
	"EXPIRE": expire,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: Ping()}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	return Value{typ: "string", str: Set(key, value)}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk
	value := Get(key)

	if value == "null" {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	return Value{typ: "string", str: Hset(hash, key, value)}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := Hget(hash, key)

	if value == "null" {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func del(args []Value) Value {
	var keys []string
	for _, arg := range args {
		keys = append(keys, arg.bulk)
	}
	value := Del(keys)

	return Value{typ: "int", num: value}

}

func expire(args []Value) Value {

	key := args[0].bulk
	numVal, err := strconv.ParseInt(args[1].bulk, 10, 64)
	if err != nil {
		return NewErrorValue(fmt.Sprintf("Invalid number for Expire:  %s ", string(args[1].bulk)))
	}
	tm := time.Duration(numVal * int64(time.Second))

	res := Expire(tm, key, nil)
	if res == -1 {
		return NewErrorValue(fmt.Sprintf("Invalid Option for Expire: option %s ", args[2].bulk))
	}
	return NewNumberValue(int16(res))
}
