package sky

type function struct {
	declaration *functionStmt
	symbols     *symbolTable
}

type returnV struct {
	value interface{}
}

func newReturnV(value interface{}) *returnV {
	return &returnV{value: value}
}

func newFunction(declaration *functionStmt, symbols *symbolTable) *function {
	return &function{
		declaration: declaration,
		symbols:     symbols,
	}
}

func (f *function) arity() int {
	return len(f.declaration.params)
}

func (f *function) call(i *Interpreter, arguments []interface{}) (ret interface{}) {
	//return value
	defer func() {
		v := recover()
		if returnValue, ok := v.(*returnV); ok {
			ret = returnValue.value
			return
		}
		panic(v)
	}()

	previous := i.currentSymbolTable
	defer func() {
		i.currentSymbolTable = previous
	}()

	//use the symbols where the function is declared, use closure environment
	i.currentSymbolTable = newSymbolTable(f.symbols)

	for k := 0; k < len(arguments); k++ {
		i.currentSymbolTable.define(f.declaration.params[k], arguments[k])
	}

	i.visitBlockStmt(f.declaration.body)
	return ret
}
