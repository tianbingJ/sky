package sky

import "fmt"

func (i *Interpreter) registerFunction() {
	i.globalSymbolTable.defineRaw("print", newPrintFunction())
}

//***********   print
type builtinPrint struct {
}

func newPrintFunction() *builtinPrint {
	return &builtinPrint{}
}

func (p *builtinPrint) arity() int {
	return 1
}

func (p *builtinPrint) call(i *Interpreter, arguments []interface{}) (ret interface{}) {
	println(fmt.Sprintf("%v", arguments[0]))
	return ret
}
