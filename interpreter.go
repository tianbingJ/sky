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
			return -iValue
		}
		if iValue, ok := value.(float64); ok {
			return -iValue
		}
		panic(newRuntimeError("Should be int valueExpr after '-' operator", op))
	}
	//bang
	if bValue, ok := value.(bool); ok {
		return !bValue
	}
	panic(newRuntimeError("Shoule be bool valueExpr after '!' operator", op))
}

func (i *interpreter) visitLiteralExpr(expression *literalExpr) interface{} {
	return expression.value
}

func (i *interpreter) visitVariableExpr(expression *variableExpr) interface{} {
	return i.current.getVariableValue(expression.token)
}

func (i *interpreter) evaluate(expression expr) interface{} {
	return expression.accept(i)
}

func (i *interpreter) visitVarStmt(varStmt *varStmt) {
	for k := 0; k < len(varStmt.elements); k++ {
		var value interface{}
		if varStmt.elements[k].valueExpr != nil {
			value = i.evaluate(varStmt.elements[k].valueExpr)
		}
		i.current.define(varStmt.elements[k].name, value)
	}
}

func (i *interpreter) visitExpressionStmt(expressionStmt *expressionStmt) {
	i.evaluate(expressionStmt.value)
}

func (i *interpreter) visitBlockStmt(block *blockStmt) {
	prev := i.current
	defer func() {
		i.current = prev
	}()
	i.current = newSymbolTable(i.current)
	for k := 0; k < len(block.statements); k++ {
		block.statements[k].accept(i)
	}
}

func (i *interpreter) visitIfStmt(ifstmt *ifStmt) {
	ifConditionValue := i.evaluate(ifstmt.ifCondition)
	if isTruthy(ifConditionValue) {
		ifstmt.ifBlock.accept(i)
		return
	}
	//elifs
	for k := 0; k < len(ifstmt.elseIfs); k++ {
		value := i.evaluate(ifstmt.elseIfs[k].condition)
		if isTruthy(value) {
			ifstmt.elseIfs[k].block.accept(i)
			return
		}
	}
	//else
	if ifstmt.elseBlock != nil {
		ifstmt.elseBlock.accept(i)
	}
}

func (i *interpreter) visitForStmt(forstmt *forStmt) {
	prev := i.current
	i.current = newSymbolTable(i.current)
	defer func() {
		i.current = prev
	}()
	defer func() {
		v := recover()
		if value, ok := v.(control_code); ok && value == break_code {
			return
		}
		panic(v)
	}()

	if forstmt.varDeclaration != nil {
		forstmt.varDeclaration.accept(i)
	} else {
		forstmt.initializers.accept(i)
	}

	for isTruthy(i.evaluate(forstmt.condition)) {
		forstmt.forBlock.accept(i)
		forstmt.increments.accept(i)
	}
}

func (i *interpreter) visitBreakStmt(breakStmt *breakStmt) {
	panic(break_code)
}

func (i *interpreter) visitWhileStmt(st *whileStmt) {
	prev := i.current
	i.current = newSymbolTable(i.current)
	defer func() {
		i.current = prev
	}()
	defer func() {
		v := recover()
		if value, ok := v.(control_code); ok && value == break_code {
			return
		}
		panic(v)
	}()
	for isTruthy(st.condition) {
		st.whileBlock.accept(i)
	}
}

func (i *interpreter) visitAssignStmt(st *assignStmt) {
	for k := 0; k < len(st.elements); k++ {
		name := st.elements[k].name
		valueExpr := st.elements[k].valueExpr
		var value interface{}
		if valueExpr != nil {
			value = i.evaluate(valueExpr)
		} else {
			value = nil
		}
		i.current.assign(name, value)
	}
}
