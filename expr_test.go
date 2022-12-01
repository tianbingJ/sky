package sky

import (
	"testing"
)

type exprTest struct {
	exprSource string
	value      interface{} ``
}

func TestArith(t *testing.T) {
	tests := []exprTest{
		{"(1+ 2) == 3\n", true},
		{"(1) + (4)\n", int64(5)},
		{"1 + 2\n", int64(3)},
		{"3 * 4 + 5\n", int64(17)},
		{" 5+ 3 * 4\n", int64(17)},
		{"1.0 + 2.0\n", 3.0},
		{"1.0 + 2\n", 3.0},
		{"1 + 1 * 2 /3.0\n", 1 + 1*2/3.0},
		{"1 << 2\n", int64(4)},
		{"1 >> 2\n", int64(0)},
		{"1024 >> 9\n", int64(2)},
		{"100 % 97\n", int64(3)},
		{"1024 | 1023\n", int64(1024 + 1023)},
		{"1024 & 1023\n", int64(0)},
		{"1024 ^ 1023\n", int64(1024 + 1023)},
		{"1 < 2\n", true},
		{"1 <= 2\n", true},
		{"1 > 2\n", false},
		{"2 >= 2\n", true},
		{"2 == 1 + 1\n", true},
		{"2 != 2\n", false},
		{"2 == 2.0\n", true},
		{"2 != 2.0\n", false},
		{"!true\n", false},
		{"!false\n", true},
		{"true and true\n", true},
		{"true and false\n", false},
		{"true or false\n", true},
		{"false or false\n", false},
		{`"hello" == "world\n"`, false},
		{`"hello" <= "world\n"`, true},
		{`"hello" > "world\n"`, false},
	}
	doTestArith(t, tests)
}

func doTestArith(t *testing.T, tests []exprTest) {
	source := joinSource(tests)
	p := getParser(source)
	exprs := make([]expr, 0)
	for !p.atEnd() {
		exprs = append(exprs, p.expressionStmt().value)
	}

	i := NewInterpreter(nil)
	results := i.interpretExpression(exprs)

	if len(results) != len(tests) {
		t.Errorf("result size not equals to test cases, result count: %d, tests count: %d", len(results), len(tests))
	}
	for idx := 0; idx < len(tests); idx++ {
		if !numberEquals(tests[idx].value, results[idx]) {
			t.Errorf("test fail for expression '%s', expect result: %v(%T), actual: %v(%T)",
				tests[idx].exprSource, tests[idx].value, tests[idx].value, results[idx], results[idx])
		}
	}
}

func joinSource(tests []exprTest) string {
	r := ""
	for i := 0; i < len(tests); i++ {
		r += tests[i].exprSource + ";"
	}
	return r
}
