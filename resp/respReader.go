package resp_basic

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type RespReader interface {
	readLine() (line []byte, n int, err error)
}

type BasicReader struct {
	buffer    []byte
	readIndex int
	lines     [][]byte
}

func NewBasicReader(rd io.Reader) *BasicReader {
	buf := make([]byte, 1024)
	length, _ := rd.Read(buf)

	return &BasicReader{buffer: buf[0:length], readIndex: 0}
}

func (r *BasicReader) linesExtractor() {
	var tmp []byte
	for _, char := range r.buffer {
		if string(char) == "\r" {
			r.lines = append(r.lines, tmp)
			fmt.Printf("Line: %v \n", string(tmp))
			tmp = nil
		} else if string(char) == "\n" {

			continue
		} else {
			tmp = append(tmp, char)

		}

	}

}

func (r *BasicReader) readLine() (line []byte, n int, err error) {
	for i := 0; i <= len(r.buffer); i++ {
		// if i == len(r.buffer) {
		// 	return line, n, errors.New("invalid end of line")
		// }
		b := r.buffer[i]

		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *BasicReader) readInteger() (x int, n int, err error) {

	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line[1]), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *BasicReader) Read() (Value, error) {
	if len(r.lines) == 0 {
		r.linesExtractor()
		fmt.Printf("lines count: %v \n", len(r.lines))
	}
	if r.readIndex == len(r.lines) {
		return Value{}, errors.New("index out of range")
	}

	_type := r.lines[r.readIndex][0]
	fmt.Printf("current line: %v \n", string(r.lines[r.readIndex]))
	fmt.Printf("current read index %v and it will be %v read \n", r.readIndex, r.readIndex+1)
	r.readIndex += 1

	switch _type {
	case ARRAY:
		fmt.Println("type array")
		return r.readArray()
	case BULK:
		fmt.Println("type bulk")
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v \n", string(_type))
		return Value{}, nil
	}

}

func (r *BasicReader) readArray() (Value, error) {
	v := Value{}
	v.Typ = "array"

	// read length of array
	arrayLen, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.Array = make([]Value, 0)
	for i := 0; i < arrayLen; i++ {

		val, err := r.Read()
		if err != nil {
			if err.Error() == "index out of range" {
				return v, nil
			}
			return v, err
		}

		// append parsed value to array
		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (r *BasicReader) readBulk() (Value, error) {
	v := Value{}

	v.Typ = "bulk"

	bulk := r.lines[r.readIndex]

	fmt.Printf("bulk value : %v \n", string(bulk))

	v.Bulk = string(bulk)
	fmt.Printf("current read index %v and it will be %v bulk \n", r.readIndex, r.readIndex+1)
	r.readIndex += 1
	return v, nil
}
