package sky

type control_code int

const (
	break_code  control_code = 1
	return_code control_code = 1
)

type stmt interface {
	accept(v stmtVisitor)
}

//*************** var stmt
type varStmt struct {
	elements []*varElement
}

type varElement struct {
	name        token
	initializer expr
}

func newVarElement(name token, initializer expr) *varElement {
	return &varElement{
		name:        name,
		initializer: initializer,
	}
}

func newVarStmt(elements []*varElement) stmt {
	return &varStmt{
		elements: elements,
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

//************* block stmt
type blockStmt struct {
	statements []stmt
}

func newBlockStmt(statements []stmt) stmt {
	return &blockStmt{
		statements: statements,
	}
}

func (block *blockStmt) accept(v stmtVisitor) {
	v.visitBlockStmt(block)
}

//*********** if stmt
type ifStmt struct {
	ifCondition expr
	ifBlock     stmt
	elseIfs     []*elseIfconditionBlock
	elseBlock   stmt
}

type elseIfconditionBlock struct {
	condition expr
	block     stmt
}

func newElseIfConditionBlock(condition expr, block stmt) *elseIfconditionBlock {
	return &elseIfconditionBlock{
		condition: condition,
		block:     block,
	}
}

func newIfStmt(ifCondition expr, ifBlock stmt, elseIfs []*elseIfconditionBlock, elseBlock stmt) stmt {
	return &ifStmt{
		ifCondition: ifCondition,
		ifBlock:     ifBlock,
		elseIfs:     elseIfs,
		elseBlock:   elseBlock,
	}
}

func (ifSt *ifStmt) accept(v stmtVisitor) {
	v.visitIfStmt(ifSt)
}

//***************** for stmt
type forStmt struct {
	// varDeclaration and initializers cannot both be set
	varDeclaration stmt
	initializers   []expr
	condition      expr
	increments     []expr
	forBlock       stmt
}

func newForStmt(varDeclaration stmt, initializers []expr, condition expr, increments []expr, forBlock stmt) stmt {
	return &forStmt{
		varDeclaration: varDeclaration,
		initializers:   initializers,
		condition:      condition,
		increments:     increments,
		forBlock:       forBlock,
	}
}

func (forst *forStmt) accept(v stmtVisitor) {
	v.visitForStmt(forst)
}

//***************

type breakStmt struct {
	breakToken token
}

func newBreakStmt(breakToken token) *breakStmt {
	return &breakStmt{
		breakToken: breakToken,
	}
}

func (breakStmt *breakStmt) accept(v stmtVisitor) {
	//TODO 分析是否在for 或者 while 内部
	v.visitBreakStmt(breakStmt)
}
