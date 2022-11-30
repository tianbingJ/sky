package sky

type control_code int

const (
	break_code control_code = 1
)

type stmt interface {
	accept(v stmtVisitor)
}

//*************** var stmt
type varStmt struct {
	elements []*assignElement
}

type assignElement struct {
	name      token
	valueExpr expr
}

func newVarElement(name token, initializer expr) *assignElement {
	return &assignElement{
		name:      name,
		valueExpr: initializer,
	}
}

func newVarStmt(elements []*assignElement) stmt {
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

func newExpressionStmt(value expr) *expressionStmt {
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

func newBlockStmt(statements []stmt) *blockStmt {
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
	initializers   *assignStmt
	condition      expr
	increments     *assignStmt
	forBlock       stmt
}

func newForStmt(varDeclaration stmt, initializers *assignStmt, condition expr, incr *assignStmt, forBlock stmt) *forStmt {
	return &forStmt{
		varDeclaration: varDeclaration,
		initializers:   initializers,
		condition:      condition,
		increments:     incr,
		forBlock:       forBlock,
	}
}

func (forst *forStmt) accept(v stmtVisitor) {
	v.visitForStmt(forst)
}

//*************** break stmt

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

//*************** while stmt
type whileStmt struct {
	condition  expr
	whileBlock stmt
}

func newWhileStmt(condition expr, whileBlock stmt) *whileStmt {
	return &whileStmt{
		condition:  condition,
		whileBlock: whileBlock,
	}
}
func (st *whileStmt) accept(v stmtVisitor) {
	//TODO 分析break是否在for 或者 while 内部
	v.visitWhileStmt(st)
}

//**************** assign stmt
type assignStmt struct {
	elements []*assignElement
}

func newAssignStmt(elements []*assignElement) *assignStmt {
	return &assignStmt{
		elements: elements,
	}
}

func (st *assignStmt) accept(v stmtVisitor) {
	v.visitAssignStmt(st)
}

//************* function stmt
type functionStmt struct {
	name   token
	params []token
	body   *blockStmt
}

func newFunctionStmt(name token, params []token, body *blockStmt) *functionStmt {
	return &functionStmt{
		name:   name,
		params: params,
		body:   body,
	}
}

func (funcStmt *functionStmt) accept(v stmtVisitor) {
	v.visitFunctionStmt(funcStmt)
}

//************** return stmt
type returnStmt struct {
	token      token
	expression expr
}

func newReturnStmt(tok token, expression expr) *returnStmt {
	return &returnStmt{
		token:      tok,
		expression: expression,
	}
}

func (retStmt *returnStmt) accept(v stmtVisitor) {
	v.visitReturnStmt(retStmt)
}
