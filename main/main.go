package main

import (
	"io/ioutil"
	"sky"
)

func main() {
	file, err := ioutil.ReadFile("./main/test.sky")
	if err != nil {
		panic(err)
	}
	lexer := sky.NewLexer(string(file))
	parser := sky.NewParser(lexer.Parse())
	stmts := parser.Parse()
	semanticResolver := sky.NewSemanticResolver()
	semanticResolver.Resolve(stmts)
	interpret := sky.NewInterpreter(semanticResolver.GetDistances())
	interpret.Interpret(stmts)
}
