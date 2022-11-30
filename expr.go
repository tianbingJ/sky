package sky

type expr interface {
	accept(v exprVisitor) interface{}
}

type binaryExpr struct {
	operator token
	left     expr
	right    expr
}

func newBinaryExpr(tokType token, left, right expr) *binaryExpr {
	return &binaryExpr{
		operator: tokType,
		left:     left,
		right:    right,
	}
}

func (l *binaryExpr) accept(v exprVisitor) interface{} {
	return v.visitBinaryExpr(l)
}

//************ liter expr

type literalExpr struct {
	value interface{}
}

func newLiteralExpr(value interface{}) *literalExpr {
	return &literalExpr{
		value: value,
	}
}

func (l *literalExpr) accept(v exprVisitor) interface{} {
	return v.visitLiteralExpr(l)
}

//************ unary expr

type unaryExpr struct {
	token      token
	expression expr
}

func newUnaryExpr(tok token, expression expr) *unaryExpr {
	return &unaryExpr{
		token:      tok,
		expression: expression,
	}
}

func (u *unaryExpr) accept(v exprVisitor) interface{} {
	return v.visitUnaryExpr(u)
}

//************ variableExpr expr

type variableExpr struct {
	token token
}

func newVariableExpr(tok token) *variableExpr {
	return &variableExpr{
		token: tok,
	}
}

func (u *variableExpr) accept(v exprVisitor) interface{} {
	return v.visitVariableExpr(u)
}

//********** call expr
type callExpr struct {
	callee    expr
	paren     token
	arguments []expr
}

func newCallExpr(callee expr, paren token, arguments []expr) *callExpr {
	return &callExpr{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}
}

func (call *callExpr) accept(v exprVisitor) interface{} {
	return v.visitCallExpr(call)
}
