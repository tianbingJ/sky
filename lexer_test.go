package sky

import (
	"testing"
)

type tokenTest struct {
	expectedTokenType tokenType
	expectedLexeme    string
	expectedLine      int
	expectedColumn    int
}

func doTest(source string, expected []tokenTest, ignoreLine, ignoreColumn bool, t *testing.T) {
	l := NewLexer(source)
	tokens := l.parse()
	expected = append(expected, tokenTest{
		expectedTokenType: EOF,
		expectedLexeme:    EOF,
		expectedLine:      0,
		expectedColumn:    0,
	})
	if len(tokens) != len(expected) {
		t.Errorf("expected name size %d, actual size %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.tokenType != expected[i].expectedTokenType {
			t.Errorf("testFail for name: %q, expected name type %s, actual type %s", tok, expected[i].expectedTokenType, tok.tokenType)
		}
		if tok.lexeme != expected[i].expectedLexeme {
			t.Errorf("testFail for name: %q, expected name lexeme %s, actual lexeme %s", tok, expected[i].expectedTokenType, tok.tokenType)
		}
		if tok.tokenType == EOF {
			continue
		}
		if !ignoreLine && tok.line != expected[i].expectedLine {
			t.Errorf("testFail for name: %q, expected name line %d, actual line %d", tok, expected[i].expectedLine, tok.line)
		}
		if !ignoreColumn && tok.column != expected[i].expectedColumn {
			t.Errorf("testFail for name: %q, expected name column %d, actual column %d", tok, expected[i].expectedColumn, tok.column)
		}
	}

}

func TestSingleOperator(t *testing.T) {
	s := `%
+-
*/`
	expected := []tokenTest{
		{MOD, MOD, 1, 1},
		{PLUS, PLUS, 2, 1},
		{MINUS, MINUS, 2, 2},
		{STAR, STAR, 3, 1},
		{SLASH, SLASH, 3, 2},
	}
	doTest(s, expected, false, false, t)
}

func TestTab(t *testing.T) {
	s := `		x	y`
	expected := []tokenTest{
		{IDENTIFIER, "x", 1, 9},
		{IDENTIFIER, "y", 1, 14},
	}
	doTest(s, expected, false, false, t)
}
func TestNextToken(t *testing.T) {
	s := `
var x = 1;
func add(a, b) {
	return a + b;
}
`
	expected := []tokenTest{
		{VAR, VAR, 2, 1},
		{IDENTIFIER, "x", 2, 5},
		{ASSIGN, "=", 2, 7},
		{INT, "1", 2, 9},
		{SEMICOLON, ";", 2, 10},
		{FUNC, FUNC, 3, 1},
		{IDENTIFIER, "add", 3, 6},
		{LPAREN, LPAREN, 3, 9},
		{IDENTIFIER, "a", 3, 10},
		{COMMA, COMMA, 3, 11},
		{IDENTIFIER, "b", 3, 13},
		{RPAREN, RPAREN, 3, 14},
		{LBRACE, LBRACE, 3, 16},
		{RETURN, RETURN, 4, 5},
		{IDENTIFIER, "a", 4, 12},
		{PLUS, PLUS, 4, 14},
		{IDENTIFIER, "b", 4, 16},
		{SEMICOLON, ";", 4, 17},
		{RBRACE, RBRACE, 5, 1},
		//{EOF, EOF},
	}
	doTest(s, expected, false, false, t)
}

func TestFor(t *testing.T) {
	s := `
//你好啊
for (var i = 1; i <= 10.2; i = i + 1) {
}`
	expected := []tokenTest{
		{FOR, FOR, 3, 1},
		{LPAREN, LPAREN, 3, 5},
		{VAR, VAR, 3, 6},
		{IDENTIFIER, "i", 3, 10},
		{ASSIGN, ASSIGN, 3, 12},
		{INT, "1", 3, 14},
		{SEMICOLON, SEMICOLON, 3, 15},
		{IDENTIFIER, "i", 3, 17},
		{LEQ, LEQ, 3, 19},
		{FLOAT, "10.2", 3, 22},
		{SEMICOLON, SEMICOLON, 3, 26},
		{IDENTIFIER, "i", 3, 28},
		{ASSIGN, ASSIGN, 3, 30},
		{IDENTIFIER, "i", 3, 32},
		{PLUS, PLUS, 3, 34},
		{INT, "1", 3, 36},
		{RPAREN, RPAREN, 3, 37},
		{LBRACE, LBRACE, 3, 39},
		{RBRACE, RBRACE, 4, 1},
	}
	doTest(s, expected, false, false, t)
}

func TestMultiComment(t *testing.T) {
	s := `
x
/**
multi line
comment
*/
y
`
	expected := []tokenTest{
		{IDENTIFIER, "x", 2, 1},
		{IDENTIFIER, "y", 7, 1},
	}
	doTest(s, expected, false, false, t)
}

func TestBitOp(t *testing.T) {
	s := `
	<< >> | ^ & and or
`
	expected := []tokenTest{
		{LSHIFT, LSHIFT, 0, 0},
		{RSHIFT, RSHIFT, 0, 0},
		{OR_BIT, OR_BIT, 0, 0},
		{XOR_BIT, XOR_BIT, 0, 0},
		{AND_BIT, AND_BIT, 0, 0},
		{AND, AND, 0, 0},
		{OR, OR, 0, 0},
	}
	doTest(s, expected, true, true, t)
}
