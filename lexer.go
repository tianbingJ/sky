package sky

import (
	"fmt"
	"strconv"
)

type lexer struct {
	source  string //source code
	current int    //next character
	line    int    //next character line
	column  int    //next character column
}

func NewLexer(source string) *lexer {
	l := &lexer{
		source:  source,
		current: 0,
		line:    1,
		column:  1,
	}
	return l
}

func (l *lexer) parse() []token {
	r := make([]token, 0)
	var tok token
	for tok = l.nextToken(); tok.tokenType != EOF; tok = l.nextToken() {
		r = append(r, tok)
	}
	r = append(r, tok)
	return r
}

func (l *lexer) nextToken() token {
startL:
	ch := l.peek()
	for ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
		l.advance()
		ch = l.peek()
		continue
	}

	switch ch {
	case 0:
		l.advance()
		return newToken(EOF, EOF, l.line, l.column-1)
	case ';':
		l.advance()
		return newToken(SEMICOLON, SEMICOLON, l.line, l.column-1)
	case '(':
		l.advance()
		return newToken(LPAREN, LPAREN, l.line, l.column-1)
	case ')':
		l.advance()
		return newToken(RPAREN, RPAREN, l.line, l.column-1)
	case '^':
		return newToken(EXP, EXP, l.line, l.column-1)
	case ',':
		l.advance()
		return newToken(COMMA, COMMA, l.line, l.column-1)
	case '+':
		l.advance()
		return newToken(PLUS, PLUS, l.line, l.column-1)
	case '-':
		l.advance()
		return newToken(MINUS, MINUS, l.line, l.column-1)
	case '*':
		l.advance()
		return newToken(STAR, STAR, l.line, l.column-1)
	case '&':
		l.advance()
		return newToken(AND_BIT, AND_BIT, l.line, l.column-1)
	case '|':
		return newToken(OR_BIT, OR_BIT, l.line, l.column-1)
	case '%':
		l.advance()
		return newToken(MOD, MOD, l.line, l.column-1)
	case '/':
		l.advance()
		if l.peek() == '/' { //singnle comment '//'
			for l.peek() != '\n' && l.peek() != 0 {
				l.advance()
			}
			goto startL
		} else if l.peek() == '*' { //multi comment  '/*'
			l.advance()
			for l.peek() != 0 {
				if l.peek() == '*' && l.peekTwo() == '/' {
					l.advance()
					l.advance()
					goto startL
				}
				l.advance()
			}
			if l.peek() != '*' || l.peekTwo() != '/' {
				panic(fmt.Sprintf("unmatched multi comment, line %d column %d", l.line, l.column))
			}
			goto startL
		}
		return newToken(SLASH, SLASH, l.line, l.column-1)
	case '{':
		l.advance()
		return newToken(LBRACE, LBRACE, l.line, l.column-1)
	case '}':
		l.advance()
		return newToken(RBRACE, RBRACE, l.line, l.column-1)
	case '!':
		l.advance()
		next := l.peek()
		if next == '=' {
			l.advance()
			return newToken(NOT_EQ, NOT_EQ, l.line, l.column-2)
		}
		return newToken(BANG, BANG, l.line, l.column-1)
	case '=':
		l.advance()
		next := l.peek()
		if next == '=' {
			l.advance()
			return newToken(EQ, EQ, l.line, l.column-2)
		}
		return newToken(ASSIGN, ASSIGN, l.line, l.column-1)
	case '<':
		l.advance()
		next := l.peek()
		if next == '=' {
			l.advance()
			return newToken(LEQ, LEQ, l.line, l.column-2)
		}
		return newToken(LT, LT, l.line, l.column-1)
	case '>':
		l.advance()
		next := l.peek()
		if next == '-' {
			l.advance()
			return newToken(GEQ, GEQ, l.line, l.column-2)
		}
		return newToken(GT, GT, l.line, l.column-1)
	case '"':
		l.advance()
		startLine := l.line
		startColumn := l.column
		startIdx := l.current
		for l.peek() != 0 && l.peek() != '"' {
			if l.peek() == '\n' {
				panic(fmt.Sprintf("string cannot cross multi-line at line %d, column %d", startLine, startColumn))
			}
			l.advance()
		}
		if l.peek() == 0 {
			panic(fmt.Sprintf("expecting '\"' at line %d  column %d ", startLine, startColumn))
		}
		strValue := l.source[startIdx:l.current]
		l.advance()
		return newToken(STRING, strValue, startLine, startColumn)
	default:
		if isDigit(ch) {
			start := l.current
			startColumn := l.column
			for isDigitElem(l.peek()) {
				l.advance()
			}
			strValue := l.source[start:l.current]
			if isInt(strValue) {
				return newToken(INT, strValue, l.line, startColumn)
			}
			if isFloat(strValue) {
				return newToken(FLOAT, strValue, l.line, startColumn)
			}
			panic(fmt.Sprintf("%s is not a number", strValue))
		} else if isLetter(ch) {
			start := l.current
			startColumn := l.column
			for isLetterELem(l.peek()) {
				l.advance()
			}
			identifier := l.source[start:l.current]
			isKeyword, tokenType_ := toKeyword(tokenType(identifier))
			if isKeyword {
				return newToken(tokenType_, identifier, l.line, startColumn)
			} else {
				return newToken(IDENTIFIER, identifier, l.line, startColumn)
			}
		}
	}
	panic(fmt.Sprintf("invalid input source at line %d column %d", l.line, l.column))
}

func isFloat(s string) bool {
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}
	return false
}

func isInt(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
func isDigitElem(ch byte) bool {
	return isDigit(ch) || ch == '.'
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isLetterELem(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || isDigit(ch) || ch == '_'
}

func (l *lexer) advance() byte {
	ch := l.peek()
	l.current++
	l.column++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else if ch == '\t' {
		l.column += 3 //add extra 3
	}
	return ch
}

func (l *lexer) peek() byte {
	if l.current >= len(l.source) {
		return 0
	}
	return l.source[l.current]
}

func (l *lexer) peekTwo() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func (l *lexer) atEnd() bool {
	return l.current == len(l.source)
}
