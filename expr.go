package sky

type visitor interface {
	visitBinaryExpr() interface{}
	visitUnaryExpr() interface{}
	visitLiteralExpr() interface{}
}

type expr interface {
	accept(v visitor) interface{}
}

type binaryExpr struct {
	operator tokenType
	left     expr
	right    expr
}

func newBinaryExpr(tokType tokenType, left, right expr) *binaryExpr {
	return &binaryExpr{
		operator: tokType,
		left:     left,
		right:    right,
	}
}

func (l *binaryExpr) accept(v visitor) interface{} {
	return v.visitBinaryExpr()
}

type literalExpr struct {
	value interface{}
}

func newLiteralExpr(value interface{}) *literalExpr {
	return &literalExpr{
		value: value,
	}
}

func (l *literalExpr) accept(v visitor) interface{} {
	return v.visitLiteralExpr()
}

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

func (u *unaryExpr) accept(v visitor) interface{} {
	return v.visitUnaryExpr()
}
