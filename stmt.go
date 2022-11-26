package sky

type stmt interface {
	accept(v stmtVisitor)
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

func (vStmt *varStmt) accept(v stmtVisitor) {
	v.visitVarStmt(vStmt)
}

//***************  expression stmt
type expressionStmt struct {
	value expr
}

func newExpressionStmt(value expr) stmt {
	return &expressionStmt{
		value: value,
	}
}

func (eStmt *expressionStmt) accept(v stmtVisitor) {
	v.visitExpressionStmt(eStmt)
}
