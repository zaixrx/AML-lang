package parser

import "aml/lexer"

// TODO: add generic return value ASAP
type StmtVisitor interface {
	VisitExpr(*ExprStmt) (Value, error);
	VisitVariableDeclaration(*VarDeclarationStmt) (Value, error);
	VisitFuncDeclarationStmt(*FuncDeclarationStmt) (Value, error);
	VisitReturn(*ReturnStmt) (Value, error);
	VisitPrint(*PrintStmt) (Value, error);
	VisitBlock(*BlockStmt) (Value, error);
	VisitBreak(*BreakStmt) (Value, error);
	VisitContinue(*ContinueStmt) (Value, error);
	VisitConditional(*ConditionalStmt) (Value, error);
	VisitWhile(*WhileStmt) (Value, error);
	VisitFor(*ForStmt) (Value, error);
}

type Stmt interface { 
	Accept(StmtVisitor) (Value, error);
}

type ExprStmt struct {
	InnerExpr Expr;
}

type VarDeclarationStmt struct {
	Name lexer.Token;
	Asset Expr;
}

type Func struct {
	Name lexer.Token;
	Params []lexer.Token;
	Body []Stmt;
}

type FuncDeclarationStmt Func;

type ReturnStmt struct {
	Asset Expr;
}

type BreakStmt struct {}

type ContinueStmt struct {}

type PrintStmt struct {
	Asset Expr;
}

type BlockStmt struct {
	Stmts []Stmt;
}

type ConditionalBranch struct {
	Condition Expr;
	NDStmt Stmt; // non-declarative statement
}

type ConditionalStmt struct {
	Branches []ConditionalBranch;
}

type WhileStmt struct {
	Cond Expr;
	NDStmt Stmt;
}

type ForStmt struct {
	Init Stmt;
	Cond Expr;
	Step Expr;
	NDStmt Stmt;
}

func (stmt *ExprStmt) Accept(vis StmtVisitor) (Value, error) {
	return vis.VisitExpr(stmt);
}

func (stmt *VarDeclarationStmt) Accept(vis StmtVisitor) (Value, error) {
	return vis.VisitVariableDeclaration(stmt);
}

func (stmt *FuncDeclarationStmt) Accept(vis StmtVisitor) (Value, error) {
	return vis.VisitFuncDeclarationStmt(stmt);
}

func (stmt *ReturnStmt) Accept(in StmtVisitor) (Value, error) {
	return in.VisitReturn(stmt);
}

func (stmt *BreakStmt) Accept(in StmtVisitor) (Value, error) {
	return in.VisitBreak(stmt);
}

func (stmt *ContinueStmt) Accept(in StmtVisitor) (Value, error) {
	return in.VisitContinue(stmt);
}

func (stmt *PrintStmt) Accept(in StmtVisitor) (Value, error) {
	return in.VisitPrint(stmt);
}

func (stmt *BlockStmt) Accept(in StmtVisitor) (Value, error) {
	return in.VisitBlock(stmt);
}

func (stmt *ConditionalStmt) Accept(in StmtVisitor) (Value, error) {
	return in.VisitConditional(stmt);
}

func (stmt *WhileStmt) Accept(in StmtVisitor) (Value, error) {
	return in.VisitWhile(stmt);
}

func (stmt *ForStmt) Accept(in StmtVisitor) (Value, error) {
	return in.VisitFor(stmt);
}
