package sky

import (
	"testing"
	"time"
)

func TestDiff(t *testing.T) {
	t1 := time.Now().UnixNano()
	sum := int64(0)
	for i := 0; i < 10000000; i++ {
		sum += int64(i)
	}
	t2 := time.Now().UnixNano()
	goTime := (t2 - t1) / 1000000
	println(goTime)

	source := `
var t1 = clock();
var sum = 0;
for var i = 0; i < 10000000; i = i + 1 {
    sum = sum + i;
}
print(sum);
var t2 = clock();
print(t2 - t1);
`
	l := NewLexer(source)
	p := NewParser(l.Parse())
	stmts := p.Parse()
	resolver := NewSemanticResolver()
	resolver.Resolve(stmts)
	i := NewInterpreter(resolver.GetDistances())
	i.Interpret(stmts)
	i.printSymbols(stmts)
	r := i.currentSymbolTable.getVariableValueRaw("sum")
	print(r)
}
