package main

import (
	"io/ioutil"
	"sky"
)

func main() {
	file, err := ioutil.ReadFile("test.sky")
	if err != nil {
		panic(err)
	}
	lexer := sky.NewLexer(string(file))
	parser := sky.NewParser(lexer.Parse())
	stmts := parser.Parse()
	interpret := sky.NewInterpreter()
	interpret.Interpret(stmts)
}
