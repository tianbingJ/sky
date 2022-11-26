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
		{"1 + 2", int64(3)},
		{"3 * 4 + 5", int64(17)},
		{" 5+ 3 * 4", int64(17)},
		{"(1) + (2)", int64(3)},
		{"1.0 + 2.0", 3.0},
		{"1.0 + 2", 3.0},
		{"1 + 1 * 2 /3.0", 1 + 1*2/3.0},
		{"1 << 2", int64(4)},
		{"1 >> 2", int64(0)},
		{"1024 >> 9", int64(2)},
		{"100 % 97", int64(3)},
		{"1024 | 1023", int64(1024 + 1023)},
		{"1024 & 1023", int64(0)},
		{"1024 ^ 1023", int64(1024 + 1023)},
		{"1 < 2", true},
		{"1 <= 2", true},
		{"1 > 2", false},
		{"2 >= 2", true},
		{"2 == 1 + 1", true},
		{"2 != 2", false},
		{"2 == 2.0", true},
		{"2 != 2.0", false},
		{`"hello" == "world"`, false},
		{`"hello" <= "world"`, true},
	}
	doTestArith(t, tests)
}

func doTestArith(t *testing.T, tests []exprTest) {
	source := joinSource(tests)
	p := getParser(source)
	exprs := make([]expr, 0)
	for !p.atEnd() {
		exprs = append(exprs, p.expression())
	}

	i := newInterpreter()
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
		r += tests[i].exprSource + " "
	}
	return r
}
