package sky

import (
	"fmt"
)

type symbolTable struct {
	symbols map[string]interface{}
	prev    *symbolTable
}

func newSymbolTable(prev *symbolTable) *symbolTable {
	return &symbolTable{
		symbols: make(map[string]interface{}),
		prev:    prev,
	}
}

func (s *symbolTable) define(name token, value interface{}) {
	if _, ok := s.symbols[name.lexeme]; ok {
		panic(newRuntimeError(fmt.Sprintf("variable %s already defined", name.lexeme), name))
	}
	s.symbols[name.lexeme] = value
}

func (s *symbolTable) getVariableValueRaw(name string) interface{} {
	symbol := s.getSymbolByVariable(name)
	if symbol == nil {
		return nil
	}
	return s.symbols[name]
}

func (s *symbolTable) getVariableValue(name token) interface{} {
	symbol := s.getSymbolByVariable(name.lexeme)
	if symbol == nil {
		panic(newRuntimeError(fmt.Sprintf("variable %s not defined", name.lexeme), name))
	}
	return symbol.symbols[name.lexeme]
}

func (s *symbolTable) assign(name token, value interface{}) {
	symbol := s.getSymbolByVariable(name.lexeme)
	if symbol == nil {
		panic(newRuntimeError(fmt.Sprintf("variable %s not defined", name.lexeme), name))
	}
	symbol.symbols[name.lexeme] = value
}

func (s *symbolTable) getSymbolByVariable(variableName string) *symbolTable {
	if _, ok := s.symbols[variableName]; ok {
		return s
	}
	if s.prev == nil {
		return nil
	}
	return s.prev.getSymbolByVariable(variableName)
}
