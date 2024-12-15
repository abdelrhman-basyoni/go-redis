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
	file    *os.File
	rd      *bufio.Reader
	mu      sync.Mutex
	rewrite bool
}

const tempFileName = "temp.aof"

var tempFile *Aof
var AOF *Aof

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}
	aof.rewrite = false
	AOF = aof
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

	if aof.rewrite {
		// if in rewrite mode store a copy to the temp file and to the temp memory
		go tempWrite(value)

	}
	fmt.Println(value.Marshal())
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

func (aof *Aof) CloseAndDelete() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Close()
	return os.Remove(aof.file.Name())
}

func DeleteAOFFile(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}
	aof.rewrite = false
	return aof, nil
}

func (aof *Aof) RewriteFile() error {
	temp, err := NewAof(tempFileName)
	fmt.Println("created  Temp AOF file")
	if err != nil {
		return fmt.Errorf("something went wrong while rewriting the AOF file")
	}
	tempFile = temp
	aof.rewrite = true
	fmt.Println("taking memory snap shot...")
	// start the rewrite process where we flush  the set and hSet  maps
	aof.mu.Lock()
	snapShot := MemDbInstance.MemorySnapShot()
	aof.mu.Unlock()

	fmt.Println("Snap shot completed successfully")
	sets := snapShot.sets.data
	hset := snapShot.hsets.data

	var newSets []goresp.Value
	var newHSets []goresp.Value

	for k, v := range sets {
		newSets = append(newSets, goresp.NewSetValue(k, v))
	}

	for h, outerV := range hset {
		for k, v := range outerV {
			newHSets = append(newHSets, goresp.NewHsetValue(h, k, v))
		}
	}
	fmt.Println("writing snap memory  to temp file...")
	bulkTempWrite(newSets)
	bulkTempWrite(newHSets)
	fmt.Println("writing snap memory  to temp file complete")
	fmt.Println("switching to new file...")
	err = aof.useTempAsMainFile()
	if err != nil {
		return fmt.Errorf("something went wrong while rewriting the AOF file %v", err)
	}
	fmt.Println("switching complete")
	aof.rewrite = false
	return nil
}
func (aof *Aof) useTempAsMainFile() error {
	tempName, mainName := tempFile.file.Name(), aof.file.Name()
	aof.mu.Lock()
	tempFile.mu.Lock()
	defer aof.mu.Unlock()
	defer tempFile.mu.Unlock()

	err := os.Rename(mainName, "toBeDeleted.aof")

	if err != nil {
		fmt.Println("something went wrong while renaming the main AOF File", err)
		return err
	}

	err = os.Rename(tempName, mainName)

	if err != nil {
		fmt.Println("something went wrong while renaming the Temp AOF File", err)
		return err
	}

	err = os.Remove("toBeDeleted.aof")
	if err != nil {
		fmt.Println("something went wrong while deleting the old AOF File", err)
		return err
	}

	err = tempFile.file.Close()
	aof.file, err = os.OpenFile(mainName, os.O_CREATE|os.O_RDWR, 0666)

	tempFile = nil
	return nil
}

func tempWrite(value goresp.Value) {
	tempFile.mu.Lock()
	defer tempFile.mu.Unlock()
	_, err := tempFile.file.Write(value.Marshal())

	if err != nil {
		fmt.Println("something went wrong while writing to the temp aof file:", err)
	}

}
func bulkTempWrite(values []goresp.Value) {
	tempFile.mu.Lock()
	defer tempFile.mu.Unlock()
	var err error
	for _, v := range values {
		_, err = tempFile.file.Write(v.Marshal())
	}

	if err != nil {
		fmt.Println("something went wrong while writing to the temp aof file:", err)
	}

}

var _, __ = NewAof("database.aof")
