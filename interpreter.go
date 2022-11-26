package sky

type interpreter struct {
	global  *symbolTable
	current *symbolTable
}

func newInterpreter() *interpreter {
	globalTable := newSymbolTable(nil)
	return &interpreter{
		global:  globalTable,
		current: globalTable,
	}
}

func (i *interpreter) interpret(statements []stmt) {
	for k := 0; k < len(statements); k++ {
		statements[k].accept(i)
	}
}

func (i *interpreter) interpretExpression(exprs []expr) []interface{} {
	r := make([]interface{}, 0)
	for idx := 0; idx < len(exprs); idx++ {
		r = append(r, exprs[idx].accept(i))
	}
	return r
}

func (i *interpreter) visitBinaryExpr(expression *binaryExpr) interface{} {
	op := expression.operator
	left := i.evaluate(expression.left)
	if op.tokenType == OR {
		//short circuit evaluation
		if isTruthy(left) {
			return left
		}
		return i.evaluate(expression.right)
	} else if op.tokenType == AND {
		if !isTruthy(left) {
			return i.evaluate(expression.left)
		}
		return i.evaluate(expression.right)
	}

	right := i.evaluate(expression.right)
	return evaluateBinary(op, left, right)
}

//! -
func (i *interpreter) visitUnaryExpr(expression *unaryExpr) interface{} {
	op := expression.token
	value := i.evaluate(expression.expression)
	if op.tokenType == MINUS {
		if iValue, ok := value.(int64); ok {
			return iValue
		}
		if iValue, ok := value.(float64); ok {
			return iValue
		}
		panic(newRuntimeError("Should be int value after '-' operator", op))
	}
	//bang
	if bValue, ok := value.(bool); ok {
		return !bValue
	}
	panic(newRuntimeError("Shoule be bool value after '!' operator", op))
}

func (i *interpreter) visitLiteralExpr(expression *literalExpr) interface{} {
	return expression.value
}

func (i *interpreter) visitVariableExpr(expression *variableExpr) interface{} {
	return i.current.getVariableValue(expression.token)
}

func (i *interpreter) visitAssignExpr(expression *assignExpr) interface{} {
	var value interface{}
	if expression.expr != nil {
		value = i.evaluate(expression.expr)
	}
	i.current.assign(expression.name, value)
	return value
}

func (i *interpreter) evaluate(expression expr) interface{} {
	return expression.accept(i)
}

func (i *interpreter) visitVarStmt(varStmt *varStmt) {
	var value interface{}
	if varStmt.initializer != nil {
		value = i.evaluate(varStmt.initializer)
	}
	i.current.define(varStmt.name, value)
}
func (i *interpreter) visitExpressionStmt(expressionStmt *expressionStmt) {
	i.evaluate(expressionStmt.value)
}
