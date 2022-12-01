package sky

type Interpreter struct {
	globalSymbolTable  *symbolTable
	currentSymbolTable *symbolTable
	distance           map[expr]int
}

func NewInterpreter(distance map[expr]int) *Interpreter {
	globalTable := newSymbolTable(nil)

	i := &Interpreter{
		globalSymbolTable:  globalTable,
		currentSymbolTable: globalTable,
		distance:           distance,
	}
	i.registerFunction()
	return i
}

func (i *Interpreter) Interpret(statements []stmt) {
	for k := 0; k < len(statements); k++ {
		statements[k].accept(i)
	}
}

func (i *Interpreter) interpretExpression(exprs []expr) []interface{} {
	r := make([]interface{}, 0)
	for idx := 0; idx < len(exprs); idx++ {
		r = append(r, exprs[idx].accept(i))
	}
	return r
}

func (i *Interpreter) visitBinaryExpr(expression *binaryExpr) interface{} {
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
func (i *Interpreter) visitUnaryExpr(expression *unaryExpr) interface{} {
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

func (i *Interpreter) visitCallExpr(expression *callExpr) interface{} {
	callee := i.evaluate(expression.callee)
	f, isFunction := callee.(callable)
	if !isFunction {
		panic(newRuntimeError("expression cannot be called", expression.paren))
	}
	arguments := make([]interface{}, 0)
	for k := 0; k < len(expression.arguments); k++ {
		arguments = append(arguments, i.evaluate(expression.arguments[k]))
	}
	return f.call(i, arguments)
}

func (i *Interpreter) call(f *function, arguments []interface{}) (ret interface{}) {
	//return value
	defer func() {
		v := recover()
		if returnValue, ok := v.(*returnV); ok {
			ret = returnValue.value
			return
		}
		panic(v)
	}()

	previous := i.currentSymbolTable
	defer func() {
		i.currentSymbolTable = previous
	}()

	//use the symbols where the function is declared, use closure environment
	i.currentSymbolTable = newSymbolTable(f.symbols)

	for k := 0; k < len(arguments); k++ {
		i.currentSymbolTable.define(f.declaration.params[k], arguments[k])
	}

	i.visitBlockStmt(f.declaration.body)
	return ret
}

func (i *Interpreter) visitLiteralExpr(expression *literalExpr) interface{} {
	return expression.value
}

func (i *Interpreter) visitVariableExpr(expression *variableExpr) interface{} {
	if v, ok := i.distance[expression]; ok {
		return i.currentSymbolTable.getVariableByDistance(v, expression.token)
	}
	return i.globalSymbolTable.getVariableValue(expression.token)
}

func (i *Interpreter) lookupVariable(name token, expression expr) interface{} {
	if dis, ok := i.distance[expression]; ok {
		return i.currentSymbolTable.getVariableByDistance(dis, name)
	}
	return i.globalSymbolTable.getVariableValue(name)
}

func (i *Interpreter) evaluate(expression expr) interface{} {
	return expression.accept(i)
}

func (i *Interpreter) visitVarStmt(varStmt *varStmt) {
	for k := 0; k < len(varStmt.elements); k++ {
		var value interface{}
		if varStmt.elements[k].valueExpr != nil {
			value = i.evaluate(varStmt.elements[k].valueExpr)
		}
		i.currentSymbolTable.define(varStmt.elements[k].name, value)
	}
}

func (i *Interpreter) visitExpressionStmt(expressionStmt *expressionStmt) {
	i.evaluate(expressionStmt.value)
}

func (i *Interpreter) visitBlockStmt(block *blockStmt) {
	prev := i.currentSymbolTable
	defer func() {
		i.currentSymbolTable = prev
	}()
	i.currentSymbolTable = newSymbolTable(i.currentSymbolTable)
	for k := 0; k < len(block.statements); k++ {
		block.statements[k].accept(i)
	}
}

func (i *Interpreter) visitIfStmt(ifstmt *ifStmt) {
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

func (i *Interpreter) visitForStmt(forstmt *forStmt) {
	prev := i.currentSymbolTable
	i.currentSymbolTable = newSymbolTable(i.currentSymbolTable)
	defer func() {
		i.currentSymbolTable = prev
	}()
	defer func() {
		v := recover()
		if value, ok := v.(control_code); ok && value == break_code {
			return
		}
		if v == nil {
			return
		}
		panic(v)
	}()

	if forstmt.varDeclaration != nil {
		forstmt.varDeclaration.accept(i)
	} else if forstmt.initializers != nil {
		forstmt.initializers.accept(i)
	}

	for isTruthy(i.evaluate(forstmt.condition)) {
		forstmt.forBlock.accept(i)
		forstmt.increments.accept(i)
	}
}

func (i *Interpreter) visitBreakStmt(breakStmt *breakStmt) {
	panic(break_code)
}

func (i *Interpreter) visitWhileStmt(st *whileStmt) {
	prev := i.currentSymbolTable
	i.currentSymbolTable = newSymbolTable(i.currentSymbolTable)
	defer func() {
		i.currentSymbolTable = prev
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

func (i *Interpreter) visitAssignStmt(st *assignStmt) {
	for k := 0; k < len(st.elements); k++ {
		name := st.elements[k].name
		valueExpr := st.elements[k].valueExpr
		var value interface{}
		if valueExpr != nil {
			value = i.evaluate(valueExpr)
		} else {
			value = nil
		}
		v, ok := i.distance[valueExpr]
		if ok {
			i.currentSymbolTable.assignByDistance(v, name, value)
		} else {
			i.globalSymbolTable.assign(name, value)
		}
	}
}

func (i *Interpreter) visitFunctionStmt(funcStmt *functionStmt) {
	f := newFunction(funcStmt, i.currentSymbolTable)
	i.currentSymbolTable.define(funcStmt.name, f)
}

func (i *Interpreter) visitReturnStmt(funcStmt *returnStmt) {
	if funcStmt.expression != nil {
		v := i.evaluate(funcStmt.expression)
		panic(newReturnV(v))
	}
	panic(newReturnV(nil))
}
