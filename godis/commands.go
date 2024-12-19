package godis

import (
	"sync"
	"time"

	goresp "github.com/abdelrhman-basyoni/goresp"
)

func Ping() string {

	return "pong"
}

type Sets struct {
	data map[string]string
	mu   sync.RWMutex
}

type Hsets struct {
	data map[string]map[string]string
	mu   sync.RWMutex
}

type MemoryDB struct {
	sets  *Sets
	hsets *Hsets
}

func NewMemoryDB() MemoryDB {
	sets := Sets{data: map[string]string{}, mu: sync.RWMutex{}}
	hsets := Hsets{data: map[string]map[string]string{}, mu: sync.RWMutex{}}

	return MemoryDB{sets: &sets, hsets: &hsets}
}

var MemDbInstance = NewMemoryDB()

func (memDb *MemoryDB) GlobalMemLock() {
	MemDbInstance.sets.mu.Lock()
	MemDbInstance.hsets.mu.Lock()
}
func (memDb *MemoryDB) GlobalMemUnlock() {
	MemDbInstance.sets.mu.Unlock()
	MemDbInstance.hsets.mu.Unlock()
}

func (memDb *MemoryDB) MemorySnapShot() MemoryDB {
	// copying sets
	newSets := &Sets{
		data: make(map[string]string),
		mu:   sync.RWMutex{},
	}
	memDb.GlobalMemLock()
	defer memDb.GlobalMemUnlock()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {

		for k, v := range memDb.sets.data {
			newSets.data[k] = v
		}
		wg.Done()
	}()

	// Create new Hsets with a deep copy of data
	newHsets := &Hsets{
		data: make(map[string]map[string]string),
		mu:   sync.RWMutex{},
	}
	go func() {

		for k, v := range memDb.hsets.data {
			// Deep copy the inner map
			newInnerMap := make(map[string]string)
			for innerK, innerV := range v {
				newInnerMap[innerK] = innerV
			}
			newHsets.data[k] = newInnerMap
		}

		wg.Done()
	}()

	wg.Wait()
	// Return the new MemoryDB
	return MemoryDB{
		sets:  newSets,
		hsets: newHsets,
	}

}

func (sets *Sets) Set(key, value string) string {
	go func() {

		if Conf.ao {
			val := goresp.NewSetValue(key, value)

			AOF.Write(val)
		}
	}()

	sets.mu.Lock()
	sets.data[key] = value
	sets.mu.Unlock()

	return "OK"
}

func (sets *Sets) Get(key string) string {

	sets.mu.RLock()
	value, ok := sets.data[key]
	sets.mu.RUnlock()

	if !ok {
		return "null"
	}

	return value
}
func (sets *Sets) Del(keys []string) int16 {
	go func() {

		if Conf.ao {
			val := goresp.NewDelValue(keys)
			AOF.Write(val)
		}
	}()

	sets.mu.RLock()
	count := int16(0)
	for _, key := range keys {
		if _, exists := sets.data[key]; exists {
			count++
			delete(sets.data, key)
		}

	}
	sets.mu.RUnlock()

	return count
}

// NX -- Set expiry only when the key has no expiry
// XX -- Set expiry only when the key has an existing expiry
// GT -- Set expiry only when the new expiry is greater than current one
// LT -- Set expiry only when the new expiry is less than current one
func (sets *Sets) Expire(expireTime time.Duration, key string, option *string) int8 {
	options := []string{"NX", "XX", "GT", "LT"}
	// TODO: handle the options
	if option != nil {
		for _, op := range options {
			if op == *option {
				break
			}

		}

		return -1
	}

	// if exists fire a goroutine that waits for the expire and then it deletes the key
	if _, exists := sets.data[key]; exists {
		go func() {
			<-time.After(expireTime)
			sets.delKey(key)
		}()

		return 1

	}

	return 0

}

// deletes a key from the Set
func (sets *Sets) delKey(key string) {
	sets.mu.Lock()
	defer sets.mu.Unlock()

	delete(sets.data, key)
}

func (hsets *Hsets) Hset(hash, key, value string) string {

	go func() {

		if Conf.ao {
			val := goresp.NewHsetValue(hash, key, value)
			AOF.Write(val)
		}
	}()

	hsets.mu.Lock()
	if _, ok := hsets.data[hash]; !ok {
		hsets.data[hash] = map[string]string{}
	}
	hsets.data[hash][key] = value
	hsets.mu.Unlock()

	return "OK"
}

func (hsets *Hsets) Hget(hash, key string) string {

	hsets.mu.RLock()
	value, ok := hsets.data[hash][key]
	hsets.mu.RUnlock()

	if !ok {
		return "null"
	}

	return value
}

func BGREWRITEAOF() error {

	return AOF.RewriteFile()

}
