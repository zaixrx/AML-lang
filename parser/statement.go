package parser

// TODO: add generic return value ASAP
type StmtVisitor interface {
	VisitExpr(*ExprStmt) error;
	VisitVariableDeclaration(*VarDeclarationStmt)  error;
	VisitPrint(*PrintStmt) error;
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

func (stmt *ExprStmt) Accept(vis StmtVisitor) error {
	return vis.VisitExpr(stmt);
}

func (stmt *VarDeclarationStmt) Accept(vis StmtVisitor) error {
	return vis.VisitVariableDeclaration(stmt);
}

func (stmt *PrintStmt) Accept(in StmtVisitor) error {
	return in.VisitPrint(stmt);
}
