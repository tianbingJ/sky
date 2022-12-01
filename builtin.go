package sky

import (
	"fmt"
	"time"
)

func (i *Interpreter) registerFunction() {
	i.globalSymbolTable.defineRaw("print", newPrintFunction())
	i.globalSymbolTable.defineRaw("clock", newClockFunction())
	i.globalSymbolTable.defineRaw("string", newStringFunction())
	i.globalSymbolTable.defineRaw("inspectSymbols", newInspectSymbols())
}

type builtin struct {
	name string
}

//***********   print
type builtinPrint struct {
	*builtin
}

func (b *builtin) String() string {
	return "func:" + b.name
}

func newPrintFunction() *builtinPrint {
	return &builtinPrint{
		&builtin{
			name: "print",
		},
	}
}

func (p *builtinPrint) arity() int {
	return 1
}

func (p *builtinPrint) call(i *Interpreter, arguments []interface{}) (ret interface{}) {
	println(fmt.Sprintf("%v", arguments[0]))
	return ret
}

//***********   clock
type clockFunction struct {
	*builtin
}

func newClockFunction() *clockFunction {
	return &clockFunction{
		&builtin{
			name: "clock",
		},
	}
}

func (c *clockFunction) arity() int {
	return 0
}

func (c *clockFunction) call(i *Interpreter, arguments []interface{}) (ret interface{}) {
	return time.Now().UnixNano() / 1000000
}

//*********** string
type stringFunction struct {
	*builtin
}

func newStringFunction() *stringFunction {
	return &stringFunction{
		&builtin{
			name: "string",
		},
	}
}

func (c *stringFunction) arity() int {
	return 1
}

func (c *stringFunction) call(i *Interpreter, arguments []interface{}) (ret interface{}) {
	return fmt.Sprintf("%v", arguments[0])
}

//*********** inspectSymbols
type inspectSymbols struct {
	*builtin
}

func newInspectSymbols() *inspectSymbols {
	return &inspectSymbols{
		&builtin{name: "inspectSymbols"},
	}
}

func (c *inspectSymbols) arity() int {
	return 1
}

func (c *inspectSymbols) call(i *Interpreter, arguments []interface{}) (ret interface{}) {
	i.printSymbols(arguments[0])
	return nil
}
