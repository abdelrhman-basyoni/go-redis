package godis

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	goresp "github.com/abdelrhman-basyoni/goresp"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
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

	if len(value.Array) == 0 {
		return false
	}
	command := strings.ToUpper(value.Array[0].Bulk)
	args := value.Array[1:]

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

	reader := goresp.NewRespReader(aof.file)

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
