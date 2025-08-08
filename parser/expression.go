package parser;

import "aml/lexer";

// TODO: add generic return value ASAP
type ExprVisitor interface {
	VisitTernary(*TernaryExpr) (Value, error);
	VisitBinary(*BinaryExpr) (Value, error);
	VisitUnary(*UnaryExpr) (Value, error);
	VisitLiteral(*LiteralExpr) (Value, error);
	VisitVariable(*VariableExpr) (Value, error);
	VisitGroup(*GroupingExpr) (Value, error);
	VisitAssign(*AssignExpr) (Value, error);
	VisitFuncCall(*FuncCall) (Value, error);
}

type Expr interface {
	Accept(ExprVisitor) (Value, error);
};

type TernaryExpr struct {
	Cond Expr;
	Iftrue Expr;
	Iffalse Expr;
};

type BinaryExpr struct {
	LOperand Expr
	Operator lexer.Token
	ROperand Expr	
};

type UnaryExpr struct {
	Operand Expr
	Operator lexer.Token
};

type LiteralExpr struct {
	ValueLiteral Value;
};

type VariableExpr struct {
	Name lexer.Token;
};

type GroupingExpr struct {
	InnerExpr Expr
};

type AssignExpr struct {
	Name string;
	Asset Expr;
};

type FuncCall struct {
	Callee Expr;
	Args[] Expr;
}

func (ter *TernaryExpr) Accept(vis ExprVisitor) (Value, error) {
	return vis.VisitTernary(ter);
}

func (bin *BinaryExpr) Accept(vis ExprVisitor) (Value, error) {
	return vis.VisitBinary(bin);
}

func (un *UnaryExpr) Accept(vis ExprVisitor) (Value, error) {
	return vis.VisitUnary(un);
}

func (lit *LiteralExpr) Accept(vis ExprVisitor) (Value, error) {
	return vis.VisitLiteral(lit);
}

func (vari *VariableExpr) Accept(vis ExprVisitor) (Value, error) {
	return vis.VisitVariable(vari);
}

func (grp *GroupingExpr) Accept(vis ExprVisitor) (Value, error) {
	return vis.VisitGroup(grp);
}

func (ass *AssignExpr) Accept(vis ExprVisitor) (Value, error) {
	return vis.VisitAssign(ass);
}

func (call *FuncCall) Accept(vis ExprVisitor) (Value, error) {
	return vis.VisitFuncCall(call);
}
