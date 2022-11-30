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
