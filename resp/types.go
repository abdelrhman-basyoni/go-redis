package resp_basic

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}
