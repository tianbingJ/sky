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
	stmts := make([]stmt, 0)
	for !p.atEnd() {
		stmts = append(stmts, p.declaration())
	}
	return stmts
}

func (p *parser) declaration() stmt {
	nextTokType := p.peek().tokenType
	if nextTokType == VAR {
		return p.varStatement()
	}
	return p.statement()
}

func (p *parser) varStatement() stmt {
	p.consumeRaw()
	name := p.consume(IDENTIFIER)
	var initializer expr
	if p.peek().tokenType == ASSIGN {
		p.consume(ASSIGN)
		initializer = p.expression()
	}
	p.consume(SEMICOLON)
	return newVarStmt(name, initializer)
}

func (p *parser) statement() stmt {
	return p.expressionStmt()
}

func (p *parser) expressionStmt() stmt {
	value := p.expression()
	p.consume(SEMICOLON)
	return newExpressionStmt(value)
}

func (p *parser) expression() expr {
	return p.assign()
}

func (p *parser) assign() expr {
	left := p.andOr()
	nextType := p.peek().tokenType
	if nextType == ASSIGN {
		name, ok := left.(*variableExpr)
		if !ok {
			panic(newSyntaxError("should be variable before '='", p.peek()))
		}
		p.consumeRaw()
		right := p.assign()
		return newAssignExpr(name.token, right)
	}
	return left
}

func (p *parser) andOr() expr {
	expression := p.bit()
	nextType := p.peek().tokenType
	for nextType == AND || nextType == OR {
		op := p.consumeRaw()
		right := p.bit()
		expression = newBinaryExpr(op, expression, right)
		nextType = p.peek().tokenType
	}
	return expression

}

func (p *parser) bit() expr {
	expression := p.eq()
	nextType := p.peek().tokenType
	for nextType == AND_BIT || nextType == OR_BIT || nextType == XOR_BIT || nextType == LSHIFT || nextType == RSHIFT {
		op := p.consumeRaw()
		right := p.eq()
		expression = newBinaryExpr(op, expression, right)
		nextType = p.peek().tokenType
	}
	return expression
}

func (p *parser) eq() expr {
	expression := p.compare()
	nextType := p.peek().tokenType
	for nextType == EQ || nextType == NOT_EQ {
		op := p.consumeRaw()
		right := p.add()
		expression = newBinaryExpr(op, expression, right)
		nextType = p.peek().tokenType
	}
	return expression
}

func (p *parser) compare() expr {
	expression := p.add()
	nextType := p.peek().tokenType
	for nextType == LEQ || nextType == LT || nextType == GT || nextType == GEQ {
		op := p.consumeRaw()
		right := p.add()
		expression = newBinaryExpr(op, expression, right)
		nextType = p.peek().tokenType
	}
	return expression
}

//+ -
func (p *parser) add() expr {
	expression := p.multiply()
	nextType := p.peek().tokenType
	for nextType == PLUS || nextType == MINUS {
		op := p.consumeRaw()
		right := p.multiply()
		expression = newBinaryExpr(op, expression, right)
		nextType = p.peek().tokenType
	}
	return expression
}

//* / %
func (p *parser) multiply() expr {
	expression := p.unary()
	nextType := p.peek().tokenType
	for nextType == STAR || nextType == SLASH || nextType == MOD {
		op := p.consumeRaw()
		rightExpr := p.unary()
		expression = newBinaryExpr(op, expression, rightExpr)
		nextType = p.peek().tokenType
	}
	return expression
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
		v, _ := strconv.ParseInt(nextTok.lexeme, 10, 64)
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
	case IDENTIFIER:
		e = newVariableExpr(nextTok)
	default:
		panic(newSyntaxError("expression", nextTok))
	}
	p.consumeRaw()
	return e
}
