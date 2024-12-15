package godis

import (
	"fmt"
	"strconv"
	"time"

	goresp "github.com/abdelrhman-basyoni/goresp"
)

type Value = goresp.Value

type commandFunction func([]Value) Value

var Handlers = map[string]commandFunction{
	"PING":         ping,
	"SET":          set,
	"GET":          get,
	"HSET":         hset,
	"HGET":         hget,
	"DEL":          del,
	"EXPIRE":       expire,
	"BGREWRITEAOF": rewriteAof,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "string", Str: Ping()}
	}

	return Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	return Value{Typ: "string", Str: MemDbInstance.sets.Set(key, value)}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk
	value := MemDbInstance.sets.Get(key)

	if value == "null" {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: value}
}

func del(args []Value) Value {
	var keys []string
	for _, arg := range args {
		keys = append(keys, arg.Bulk)
	}
	value := MemDbInstance.sets.Del(keys)

	return Value{Typ: "int", Num: value}

}

func expire(args []Value) Value {

	key := args[0].Bulk
	numVal, err := strconv.ParseInt(args[1].Bulk, 10, 64)
	if err != nil {
		return goresp.NewErrorValue(fmt.Sprintf("Invalid number for Expire:  %s ", string(args[1].Bulk)))
	}
	tm := time.Duration(numVal * int64(time.Second))

	res := MemDbInstance.sets.Expire(tm, key, nil)
	if res == -1 {
		return goresp.NewErrorValue(fmt.Sprintf("Invalid Option for Expire: option %s ", args[2].Bulk))
	}
	return goresp.NewNumberValue(int16(res))
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	return Value{Typ: "string", Str: MemDbInstance.hsets.Hset(hash, key, value)}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := MemDbInstance.hsets.Hget(hash, key)

	if value == "null" {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: value}
}

func rewriteAof(args []Value) Value {

	err := BGREWRITEAOF()
	if err != nil {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	return Value{Typ: "string", Str: "OK"}

}
