package sky

type stmt interface {
	accept(v stmtVisitor) interface{}
}

//*************** var stmt
type varStmt struct {
	name        token
	initializer expr
}

func newVarStmt(name token, initializer expr) stmt {
	return &varStmt{
		name:        name,
		initializer: initializer,
	}
}

func (vStmt *varStmt) accept(v stmtVisitor) interface{} {
	return v.visitVarStmt()
}

//***************  expression stmt
type expressionStmt struct {
	name  token
	value expr
}

func newExpressionStmt(value expr) stmt {
	return &expressionStmt{
		value: value,
	}
}

func (eStmt *expressionStmt) accept(v stmtVisitor) interface{} {
	return v.visitVarStmt()
}
