package godis

import "sync"

func Ping() string {

	return "pong"
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func Set(key, value string) string {

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return "OK"
}

func Get(key string) string {

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return "null"
	}

	return value
}

func Hset(hash, key, value string) string {

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return "OK"
}

func Hget(hash, key string) string {

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return "null"
	}

	return value
}
