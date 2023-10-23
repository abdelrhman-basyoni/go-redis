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
			// buf := make([]byte, 1024)
			// length, _ := buff.Read(buf)
			// fmt.Println(string(buf[:length]))
			return NewRespIo(buff), nil
		}
	default:
		return nil, errors.New("invalid buffer type")
	}
}
