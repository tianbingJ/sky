package sky

import "strconv"

type parser struct {
	tokens  []token
	current int //指向下一个token
}

func newParser(tokens []token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) atEnd() bool {
	if p.current >= len(p.tokens) {
		return true
	}
	return p.tokens[p.current].tokenType == EOF
}

//just see next token
func (p *parser) check(tokType tokenType) bool {
	if p.atEnd() {
		return false
	}
	return p.peek().tokenType == tokType
}

//just move to next token
func (p *parser) consumeRaw() token {
	tok := p.peek()
	p.current++
	return tok
}

func (p *parser) peek() token {
	return p.tokens[p.current]
}

func (p *parser) previous() token {
	return p.tokens[p.current-1]
}

//check and move to next token
func (p *parser) consume(tokType tokenType) token {
	if !p.check(tokType) {
		panic(newTokenSyntaxError(tokType, p.peek()))
	}
	p.current++
	return p.previous()
}

func (p *parser) parse() []stmt {
	return nil
}

func (p *parser) unary() expr {
	if p.check(BANG) || p.check(MINUS) {
		tok := p.peek()
		p.consumeRaw()
		expression := p.unary()
		return newUnaryExpr(tok, expression)
	}
	return p.primary()
}

func (p *parser) primary() expr {
	nextTok := p.peek()
	var e expr = nil
	switch nextTok.tokenType {
	case INT:
		v, _ := strconv.Atoi(nextTok.lexeme)
		e = newLiteralExpr(v)
	case FLOAT:
		v, _ := strconv.ParseFloat(nextTok.lexeme, 64)
		e = newLiteralExpr(v)
	case STRING:
		e = newLiteralExpr(nextTok.lexeme)
	case TRUE:
		e = newLiteralExpr(true)
	case FALSE:
		e = newLiteralExpr(false)
	case NIL:
		e = newLiteralExpr(nil)
	default:
		panic(newSyntaxError("expressoin", nextTok))
	}
	p.consumeRaw()
	return e
}
