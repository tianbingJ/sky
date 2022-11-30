package sky

import (
	"fmt"
	"testing"
)

func TestVarStmt(t *testing.T) {
	source := `var a, x = 1, "string";`
	p := getParser(source)
	statement := p.Parse()
	in := NewInterpreter()
	in.Interpret(statement)
	valueA := in.currentSymbolTable.getVariableValueRaw("a")
	valueX := in.currentSymbolTable.getVariableValueRaw("x")
	if !numberEquals(valueA, int64(1)) {
		t.Errorf("var not assign to varible effective, valueExpr should be 1, got %v", valueA)
	}
	if valueX != "string" {
		t.Errorf("var not assign to varible effective, valueExpr should be 'string', got %v", valueX)
	}
}

func TestBlockStmt(t *testing.T) {
	source := `
var a = 1; 
var r1;
var r2;
var r3;
r1 = a;
{
var a = 2;
r2 = a;
}
r3 = a;
`
	p := getParser(source)
	statement := p.Parse()
	in := NewInterpreter()
	in.Interpret(statement)
	value1 := in.currentSymbolTable.getVariableValueRaw("r1")
	value2 := in.currentSymbolTable.getVariableValueRaw("r2")
	value3 := in.currentSymbolTable.getVariableValueRaw("r3")
	if !numberEquals(value1, int64(1)) {
		t.Errorf("var not assign to varible effective, valueExpr should be 1, got %v", value1)
	}
	if !numberEquals(value2, int64(2)) {
		t.Errorf("var not assign to varible effective, valueExpr should be 2, got %v", value2)
	}
	if !numberEquals(value3, int64(1)) {
		t.Errorf("var not assign to varible effective, valueExpr should be 1, got %v", value3)
	}
}

func TestIfStmt(t *testing.T) {
	values := []int{1, 5, 10, 0}
	result := []int{11, 55, 100, -1}
	for i := 0; i < len(values); i++ {
		source := fmt.Sprintf(`
var x = %d;
var r;
if x == 1 {
	r = 11;
} elif x == 5{
	r = 55;
} elif x == 10 {
	r = 100;
} else {
	r = -1;
}
`, values[i])
		p := getParser(source)
		statement := p.Parse()
		in := NewInterpreter()
		in.Interpret(statement)
		r := in.currentSymbolTable.getVariableValueRaw("r")
		if r != int64(result[i]) {
			t.Errorf("expect 'r' = %d, but got %d", result[i], r)
		}
	}
}

func TestForStmt(t *testing.T) {
	source := `
var sum = 0;
for var i , j = 0, 0; i < 12; i , j = i + 1 , j + 1 {
	if i >= 10 {
		break;
	}
	sum = sum + i + j;
}
`
	p := getParser(source)
	statement := p.Parse()
	in := NewInterpreter()
	in.Interpret(statement)
	sum := in.currentSymbolTable.getVariableValueRaw("sum")
	if sum != int64(90) {
		t.Errorf("expect 'sum' = %d, but got %d", 90, sum)
	}
}

func TestWhileStmt(t *testing.T) {
	source := `
var sum , i = 0, 0;
while i < 20 {
	if i >= 10 {
		break;
	}
	sum = sum + i;
	i = i + 1;
}
`
	p := getParser(source)
	statement := p.Parse()
	in := NewInterpreter()
	in.Interpret(statement)
	sum := in.currentSymbolTable.getVariableValueRaw("sum")
	if sum != int64(45) {
		t.Errorf("expect 'sum' = %d, but got %d", 45, sum)
	}
}

func TestFunction(t *testing.T) {
	source := `
var sum = 0;
func add(x, y) {
	return x + y;
}
var f = add;
sum = f(1 , 2);
`
	p := getParser(source)
	statement := p.Parse()
	in := NewInterpreter()
	in.Interpret(statement)
	sum := in.currentSymbolTable.getVariableValueRaw("sum")
	if sum != int64(3) {
		t.Errorf("expect 'sum' = %d, but got %d", 3, sum)
	}
}
