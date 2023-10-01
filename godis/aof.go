package godis

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Aof struct {
	file  *os.File
	rd    *bufio.Reader
	mu    sync.Mutex
	respW RespWriter
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func parser(value Value) bool {

	if len(value.array) == 0 {
		return false
	}
	command := strings.ToUpper(value.array[0].bulk)
	args := value.array[1:]

	handler, ok := Handlers[command]
	if !ok {
		fmt.Println("Invalid command: ", command)
		return false
	}

	handler(args)

	return true
}

func (aof *Aof) Read() error {

	aof.file.Seek(0, io.SeekStart)

	reader, _ := NewRespReader(aof.file)

	for {
		value, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
		parser(value)

	}
	fmt.Println("AOF Finished recovering")
	return nil
}

var AOF, _ = NewAof("database.aof")
