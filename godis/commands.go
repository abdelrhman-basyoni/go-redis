package godis

import (
	"sync"
	"time"
)

func Ping() string {

	return "pong"
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func Set(key, value string) string {
	go func() {

		if Conf.ao {
			val := NewSetValue(key, value)
			AOF.Write(val)
		}
	}()

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

	go func() {

		if Conf.ao {
			val := NewHsetValue(hash, key, value)
			AOF.Write(val)
		}
	}()

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

func Del(keys []string) int16 {
	go func() {

		if Conf.ao {
			val := NewDelValue(keys)
			AOF.Write(val)
		}
	}()

	SETsMu.RLock()
	count := int16(0)
	for _, key := range keys {
		if _, exists := SETs[key]; exists {
			count++
			delete(SETs, key)
		}

	}
	SETsMu.RUnlock()

	return count
}

// NX -- Set expiry only when the key has no expiry
// XX -- Set expiry only when the key has an existing expiry
// GT -- Set expiry only when the new expiry is greater than current one
// LT -- Set expiry only when the new expiry is less than current one
func Expire(expireTime time.Duration, key string, option *string) int8 {
	options := []string{"NX", "XX", "GT", "LT"}
	if option != nil {
		for _, op := range options {
			if op == *option {
				break
			}

		}

		return -1
	}

	// if exists fire a goroutine that waits for the expire and then it deletes the key
	if _, exists := SETs[key]; exists {
		go func() {
			<-time.After(expireTime)
			delKey(key)
		}()

		return 1

	}

	return 0

}

// deletes a key from the Set
func delKey(key string) {
	SETsMu.Lock()
	defer SETsMu.Unlock()

	delete(SETs, key)
}
