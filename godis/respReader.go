package godis

import (
	"errors"
	"io"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type RespReader interface {
	// readLine() (line []byte, n int, err error)
	Read() (Value, error)
}

func NewRespReader(buff interface{}) (RespReader, error) {
	switch buff := buff.(type) {
	case string:
		{
			return NewBasicReader(buff), nil

		}
	case io.Reader:
		{
			return NewRespIo(buff), nil
		}
	default:
		return nil, errors.New("invalid buffer type")
	}
}
