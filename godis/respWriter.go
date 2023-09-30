package godis

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

func (w *BasicWriter) Write(bytes []byte) error {

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
