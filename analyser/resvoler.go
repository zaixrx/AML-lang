package analyzer

import (
	"aml/parser"
	"aml/interpreter"
)

type Resolver struct {
	in interpreter.Interpreter;
}

type Value = parser.Value;

func (res *Resolver) declare(name string) {
}
func (res *Resolver) resolve(name string);
func (res *Resolver) define(name string);

func (ana *Resolver) VisitVariableDeclaration(*parser.VarDeclarationStmt) (Value, error);
func (ana *Resolver) VisitFuncDeclarationStmt(*parser.FuncDeclarationStmt) (Value, error);
func (ana *Resolver) VisitVariable(*parser.VariableExpr) (Value, error);
func (ana *Resolver) VisitFuncCall(*parser.FuncCall) (Value, error);
func (ana *Resolver) VisitAssign(*parser.AssignExpr) (Value, error);
func (ana *Resolver) VisitLiteral(*parser.LiteralExpr) (Value, error);
func (ana *Resolver) VisitUnary(*parser.UnaryExpr) (Value, error);
func (ana *Resolver) VisitBinary(*parser.BinaryExpr) (Value, error);
func (ana *Resolver) VisitTernary(*parser.TernaryExpr) (Value, error);
func (ana *Resolver) VisitGroup(*parser.GroupingExpr) (Value, error);
func (ana *Resolver) VisitExpr(*parser.ExprStmt) (Value, error);
func (ana *Resolver) VisitReturn(*parser.ReturnStmt) (Value, error);
func (ana *Resolver) VisitBreak(*parser.BreakStmt) (Value, error);
func (ana *Resolver) VisitContinue(*parser.ContinueStmt) (Value, error);
func (ana *Resolver) VisitConditional(*parser.ConditionalStmt) (Value, error);
func (ana *Resolver) VisitWhile(*parser.WhileStmt) (Value, error);
func (ana *Resolver) VisitFor(*parser.ForStmt) (Value, error);
func (ana *Resolver) VisitPrint(*parser.PrintStmt) (Value, error);
func (ana *Resolver) VisitBlock(*parser.BlockStmt) (Value, error);

func NewResolver(in interpreter.Interpreter) *Resolver {
	return &Resolver{
		in: in,
	};
}

func (ana *Resolver) Resolve(stmt parser.Stmt) (Value, error) {
	return stmt.Accept(ana);
}
