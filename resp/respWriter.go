package resp_basic

import "io"

type RespWriter interface {
	Write(v Value) error
}

type BasicWriter struct {
	writer io.Writer
}

func NewBasicWriter(w io.Writer) *BasicWriter {
	return &BasicWriter{writer: w}
}

func (w *BasicWriter) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (w *BasicWriter) newSetValue(key, value string) Value {

	arr := []Value{Value{Typ: "bulk", Bulk: key}, Value{Typ: "bulk", Bulk: value}}

	return Value{Typ: "array", Array: arr}
}

func (w *BasicWriter) newHSetValue(key, value, hash string) Value {

	arr := []Value{Value{Typ: "bulk", Bulk: hash}, Value{Typ: "bulk", Bulk: key}, Value{Typ: "bulk", Bulk: value}}

	return Value{Typ: "array", Array: arr}
}
