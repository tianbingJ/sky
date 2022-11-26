package sky

import "testing"

func TestVarStmt(t *testing.T) {
	source := `var a = 1; var x = "string";`
	p := getParser(source)
	statement := p.parse()
	in := newInterpreter()
	in.interpret(statement)
	valueA := in.current.getVariableValueRaw("a")
	valueX := in.current.getVariableValueRaw("x")
	if !numberEquals(valueA, int64(1)) {
		t.Errorf("var not assign to varible effective, value should be 1, got %v", valueA)
	}
	if valueX != "string" {
		t.Errorf("var not assign to varible effective, value should be 'string', got %v", valueX)
	}
}
