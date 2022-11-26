package sky

type interpreter struct {
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
	//TODO
	return nil
}

func (i *interpreter) visitAssignExpr(expression *assignExpr) interface{} {
	//TODO
	return nil
}

func (i *interpreter) evaluate(expression expr) interface{} {
	return expression.accept(i)
}
