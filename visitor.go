package sky

type exprVisitor interface {
	visitBinaryExpr(expression *binaryExpr) interface{}
	visitUnaryExpr(expression *unaryExpr) interface{}
	visitLiteralExpr(expression *literalExpr) interface{}
	visitVariableExpr(expression *variableExpr) interface{}
	visitAssignExpr(expression *assignExpr) interface{}
}

type stmtVisitor interface {
	visitVarStmt(varStmt *varStmt)
	visitExpressionStmt(expressionStmt *expressionStmt)
}
