package sky

import "fmt"

type syntaxTokenError struct {
	expected tokenType
	actual   token
}

func (s syntaxTokenError) Error() string {
	return fmt.Sprintf("Unexpected name %s, Expecting %s at line %d  column %d\n",
		s.actual.tokenType, s.expected, s.actual.line, s.actual.column)
}

func newTokenSyntaxError(expected tokenType, tok token) error {
	return syntaxTokenError{
		expected: expected,
		actual:   tok,
	}
}

type syntaxError struct {
	expected string
	token    token
}

func (s syntaxError) Error() string {
	return fmt.Sprintf("Unexpected name %s, Expecting %s at line %d  column %d\n",
		s.token.tokenType, s.expected, s.token.line, s.token.column)
}

func newSyntaxError(expected string, tok token) error {
	return syntaxError{
		expected: expected,
		token:    tok,
	}
}
