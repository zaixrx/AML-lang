package parser

// TODO: add generic return value ASAP
type StmtVisitor interface {
	VisitExpr(*ExprStmt) error;
	VisitVariableDeclaration(*VarDeclarationStmt)  error;
	VisitPrint(*PrintStmt) error;
	VisitBlock(*BlockStmt) error;
}

type Stmt interface { 
	Accept(StmtVisitor) error;
}

type ExprStmt struct {
	InnerExpr Expr;
}

type VarDeclarationStmt struct {
	Name string;
	Asset Expr;
}

type PrintStmt struct {
	Asset Expr;
}

type BlockStmt struct {
	Stmts []Stmt;	
}

type Statmetsd struct {};

func (stmt *ExprStmt) Accept(vis StmtVisitor) error {
	return vis.VisitExpr(stmt);
}

func (stmt *VarDeclarationStmt) Accept(vis StmtVisitor) error {
	return vis.VisitVariableDeclaration(stmt);
}

func (stmt *PrintStmt) Accept(in StmtVisitor) error {
	return in.VisitPrint(stmt);
}

func (stmt *BlockStmt) Accept(in StmtVisitor) error {
	return in.VisitBlock(stmt);
}
