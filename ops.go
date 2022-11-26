package sky

func evaluateBinary(operator token, left, right interface{}) interface{} {
	if operator.tokenType == PLUS {
		ls, okl := left.(string)
		rs, okr := right.(string)
		if okl || okr {
			if !okl {
				panic(newRuntimeError(operator.lexeme+" left is a not a string", operator))
			}
			if !okr {
				panic(newRuntimeError(operator.lexeme+" right is a not a string", operator))
			}
			return ls + rs
		}
	}

	switch operator.tokenType {
	case OR_BIT, AND_BIT, XOR_BIT, LSHIFT, RSHIFT, MOD:
		ilvalue, okl := left.(int64)
		irvalue, okr := right.(int64)
		if !okl {
			panic(newRuntimeError(operator.lexeme+" left expression should be int", operator))
		}
		if !okr {
			panic(newRuntimeError(operator.lexeme+" right expression should be int", operator))
		}
		switch operator.tokenType {
		case OR_BIT:
			return ilvalue | irvalue
		case AND_BIT:
			return ilvalue & irvalue
		case XOR_BIT:
			return ilvalue ^ irvalue
		case MOD:
			return ilvalue % irvalue
		}

		//shift operator right expression should >= 0
		if irvalue < 0 {
			panic(newRuntimeError("right value should not be negative", operator))
		}
		switch operator.tokenType {
		case LSHIFT:
			return ilvalue << irvalue
		case RSHIFT:
			return ilvalue >> irvalue
		}
	case STAR, SLASH, PLUS, MINUS:
		checkBothNumbers(left, right, operator)
		if hasFloat(left, right) {
			fl := toFloat(left, operator)
			fr := toFloat(right, operator)
			switch operator.tokenType {
			case STAR:
				return fl * fr
			case SLASH:
				return fl / fr
			case PLUS:
				return fl + fr
			case MINUS:
				return fl - fr
			}
		}
		il, _ := left.(int64)
		ir, _ := right.(int64)
		switch operator.tokenType {
		case STAR:
			return il * ir
		case SLASH:
			return il / ir
		case PLUS:
			return il + ir
		case MINUS:
			return il - ir
		}
	default:
	}
	panic(newRuntimeError("unsupported operator", operator))
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if bvalue, ok := value.(bool); ok {
		return bvalue
	}
	return true
}

func toFloat(value interface{}, operator token) float64 {
	if v, ok := value.(float64); ok {
		return v
	}
	if v, ok := value.(int64); ok {
		return float64(v)
	}
	panic(newRuntimeError("not reachable here", operator))
}

func hasFloat(left, right interface{}) bool {
	_, leftFloat := left.(float64)
	_, rightFlot := right.(float64)
	return leftFloat || rightFlot
}

func checkBothNumbers(left, right interface{}, operator token) {
	_, lInt := left.(int64)
	_, rInt := right.(int64)
	_, lFloat := left.(float64)
	_, rFloat := right.(int64)
	bothNumbers := (lInt || lFloat) && (rInt || rFloat)
	if !bothNumbers {
		panic(newRuntimeError("both value should be numbers", operator))
	}
}
