package sky

import "fmt"

func (i *Interpreter) printSymbols(prompt interface{}) {
	fmt.Printf("---------%v--------\n", prompt)
	i.printSymbolTable(i.currentSymbolTable)
}

func (i *Interpreter) printSymbolTable(current *symbolTable) int {
	if current == nil {
		return 0
	}
	depth := 0
	if current.prev != nil {
		depth = i.printSymbolTable(current.prev)
	}

	for key, element := range current.symbols {
		fmt.Printf(printLeadingTab(depth)+"[%s, %v]\n", key, element)
	}
	return depth + 1
}

func printLeadingTab(times int) string {
	r := ""
	for i := 0; i < times; i++ {
		r += "\t"
	}
	return r
}
