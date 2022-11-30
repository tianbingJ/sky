package sky

type callable interface {
	arity() int

	call(i *Interpreter, arguments []interface{}) (ret interface{})
}
