package godis

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}
