package sky

import "fmt"

type semanticResolver struct {
	inFunction bool
	inFor      bool
	inWhile    bool
	scopes     []*scope
	distances  map[expr]int
}

type varState int

const (
	DECLARED varState = 1
	DEFINED  varState = 2
)

type scope struct {
	symbols map[string]varState
}

func newScope() *scope {
	return &scope{
		symbols: make(map[string]varState),
	}
}

func NewSemanticResolver() *semanticResolver {
	return &semanticResolver{
		inFor:      false,
		inWhile:    false,
		inFunction: false,
		scopes:     make([]*scope, 0),
	}
}

func (s *semanticResolver) Resolve(stmts []stmt) {
	for i := 0; i < len(stmts); i++ {
		stmts[i].accept(s)
	}
}

func (s *semanticResolver) peekScope() *scope {
	return s.scopes[len(s.scopes)-1]
}

//just consider locals
func (s *semanticResolver) declare(name token) {
	if len(s.scopes) == 0 {
		return
	}
	tscope := s.peekScope()
	if _, ok := tscope.symbols[name.lexeme]; ok {
		panic(newSyntaxError(fmt.Sprintf("Variable %s already declared in current scope", name.lexeme), name))
	}
	tscope.symbols[name.lexeme] = DECLARED
}

//just consider locals
func (s *semanticResolver) define(name token) {
	if len(s.scopes) == 0 {
		return
	}
	tscope := s.peekScope()
	tscope.symbols[name.lexeme] = DEFINED
}

func (s *semanticResolver) visitVarStmt(varStmt *varStmt) {
	//declare
	for i := 0; i < len(varStmt.elements); i++ {
		element := varStmt.elements[i]
		s.declare(element.name)
	}

	//visitVarStmt
	for i := 0; i < len(varStmt.elements); i++ {
		element := varStmt.elements[i]
		if element.valueExpr != nil {
			s.resolveExpression(element.valueExpr)
		}
	}
	//define
	for i := 0; i < len(varStmt.elements); i++ {
		element := varStmt.elements[i]
		s.define(element.name)
	}
}

func (s *semanticResolver) visitExpressionStmt(expressionStmt *expressionStmt) {
	s.resolveExpression(expressionStmt.value)
}

func (s *semanticResolver) visitBlockStmt(block *blockStmt) {
	s.beginScope()
	defer func() {
		s.endScope()
	}()
	for _, statement := range block.statements {
		s.resolveStmt(statement)
	}
}

func (s *semanticResolver) visitIfStmt(ifstmt *ifStmt) {
	s.resolveExpression(ifstmt.ifCondition)
	s.resolveStmt(ifstmt.ifBlock)

	if len(ifstmt.elseIfs) > 0 {
		for _, v := range ifstmt.elseIfs {
			s.resolveExpression(v.condition)
			s.resolveStmt(v.block)
		}
	}
	if ifstmt.elseBlock != nil {
		s.resolveStmt(ifstmt.elseBlock)
	}
}

func (s *semanticResolver) visitForStmt(forstmt *forStmt) {
	inFor := s.inFor
	s.inFor = true
	defer func() {
		s.inFor = inFor
	}()
	if forstmt.varDeclaration != nil {
		s.resolveStmt(forstmt.varDeclaration)
	}
	if forstmt.initializers != nil {
		s.resolveStmt(forstmt.initializers)
	}
	if forstmt.condition != nil {
		s.resolveExpression(forstmt.condition)
	}
	s.resolveStmt(forstmt.forBlock)
	if forstmt.increments != nil {
		s.resolveStmt(forstmt.increments)
	}
}

func (s *semanticResolver) visitBreakStmt(breakStmt *breakStmt) {
	if !s.inFor && !s.inWhile {
		panic(newSyntaxError("'break' must be used in 'for' or 'while' stmt.", breakStmt.breakToken))
	}
}

func (s *semanticResolver) visitWhileStmt(whileStmt *whileStmt) {
	inWhile := s.inWhile
	s.inWhile = true
	defer func() {
		s.inWhile = inWhile
	}()
	s.resolveExpression(whileStmt.condition)
	s.resolveStmt(whileStmt.whileBlock)
}

func (s *semanticResolver) visitAssignStmt(assign *assignStmt) {
	for _, elem := range assign.elements {
		s.resolveExpression(elem.valueExpr)
		s.resolve(elem.valueExpr, elem.name)
	}
}

func (s *semanticResolver) resolve(expression expr, name token) {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		if _, ok := s.scopes[i].symbols[name.lexeme]; ok {
			s.distances[expression] = len(s.scopes) - 1 - i
		}
	}
}

func (s *semanticResolver) beginScope() {
	s.scopes = append(s.scopes, newScope())
}

func (s *semanticResolver) endScope() {
	s.scopes = s.scopes[0 : len(s.scopes)-1]
}

func (s *semanticResolver) visitFunctionStmt(funcStmt *functionStmt) {
	s.declare(funcStmt.name)
	inFunction := s.inFunction

	s.beginScope()
	defer func() {
		s.inFunction = inFunction
		s.endScope()
		s.define(funcStmt.name)
	}()
	s.inFunction = true
	for _, v := range funcStmt.params {
		s.declare(v)
		s.define(v)
	}
	s.resolveStmt(funcStmt.body)
}

func (s *semanticResolver) visitReturnStmt(stmt *returnStmt) {
	if !s.inFunction {
		panic(newSyntaxError("'return' must be in 'function' stmt.", stmt.token))
	}
	if stmt.expression != nil {
		s.resolveExpression(stmt.expression)
	}
}

func (s *semanticResolver) visitBinaryExpr(expression *binaryExpr) interface{} {
	s.resolveExpression(expression.left)
	s.resolveExpression(expression.right)
	return nil
}

func (s *semanticResolver) visitUnaryExpr(expression *unaryExpr) interface{} {
	s.resolveExpression(expression.expression)
	return nil
}

func (s *semanticResolver) visitLiteralExpr(expression *literalExpr) interface{} {
	return nil
}

func (s *semanticResolver) visitVariableExpr(expression *variableExpr) interface{} {
	if len(s.scopes) > 0 {
		if v, ok := s.peekScope().symbols[expression.token.lexeme]; ok && v == DECLARED {
			panic(newSyntaxError(fmt.Sprintf("%s is not defined", expression.token.lexeme), expression.token))
		}
	}
	return nil
}

func (s *semanticResolver) visitCallExpr(expression *callExpr) interface{} {
	s.resolveExpression(expression.callee)
	for i := 0; i < len(expression.arguments); i++ {
		s.resolveExpression(expression.arguments[i])
	}
	return nil
}

func (s *semanticResolver) resolveExpression(expression expr) {
	expression.accept(s)
}

func (s *semanticResolver) resolveStmt(statement stmt) {
	statement.accept(s)
}
