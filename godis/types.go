package godis

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

type RespWriter interface {
	Write(v Value) error
}
