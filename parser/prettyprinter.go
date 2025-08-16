package parser

import (
	"aml/lexer"
	"fmt"
	"strings"
)

type PrettyPrinter struct {
	tabs int;
	builder strings.Builder;
}

func (p *PrettyPrinter) tab() {
	p.tabs++;
}

func (p *PrettyPrinter) untab() {
	p.tabs--;
}

func (p *PrettyPrinter) append_raw(strs ...string) {
	p.builder.WriteString(strings.Join(strs, ""));
}

func (p *PrettyPrinter) append(strs ...string) {
	p.builder.WriteString(strings.Repeat("  ", p.tabs));
	p.builder.WriteString(strings.Join(strs, ""));
}

func (p *PrettyPrinter) appendln(strs ...string) {
	strs = append(strs, "\n");
	p.append(strs...);
}

func (p *PrettyPrinter) print_header(header string) {
	p.appendln("(", header, ")");
}

func (p *PrettyPrinter) print_def_value(name string, vals ...Value) {
	p.append(name, ": [");
	for i, val := range vals {
		p.append_raw(fmt.Sprint(val));
		if i + 1 != len(vals) {
			p.append_raw(", ");
		}
	}
	p.append_raw("]\n");
}

func (p *PrettyPrinter) print_def_token(name string, toks ...lexer.Token) {
	p.append(name, ": [");
	for i, tok := range toks {
		p.append_raw(string(tok.Lexeme));
		if i + 1 != len(toks) {
			p.append_raw(", ");
		}
	}
	p.append_raw("]\n");
}

func (p *PrettyPrinter) print_def_expr(name string, exprs ...Expr) {
	p.append(name, ":");
	p.tab();
		for _, expr := range exprs {
			p.append("\n");
			expr.Accept(p);
		}
	p.untab();
}

func (p *PrettyPrinter) print_def_stmt(name string, stmts ...Stmt) {
	p.append(name, ":");
	p.tab();
		for _, stmt := range stmts {
			p.append("\n");
			stmt.Accept(p);
		}
	p.untab();
}

func (p *PrettyPrinter) VisitTernary(ter *TernaryExpr) (Value, error) {
	p.print_header("Ternary");
	p.tab();
		p.print_def_expr("Cond", ter.Cond);
		p.print_def_expr("IfTrue", ter.Iftrue);
		p.print_def_expr("IfFalse", ter.Iffalse);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitBinary(bin *BinaryExpr) (Value, error) {
	p.print_header("Binary");
	p.tab();
		p.print_def_expr("LOperand", bin.LOperand);
		p.print_def_token("Operator", bin.Operator);
		p.print_def_expr("ROperand", bin.ROperand);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitUnary(un *UnaryExpr) (Value, error) {
	p.print_header("Unary");
	p.tab();
		p.print_def_expr("Operand", un.Operand);
		p.print_def_token("Operator", un.Operator);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitLiteral(lit *LiteralExpr) (Value, error) {
	p.print_header("Literal");
	p.tab();
		p.print_def_value("Value", lit.ValueLiteral);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitVariable(vari *VariableExpr) (Value, error) {
	p.print_header("Variable");
	p.tab();
		p.print_def_token("Name", vari.Name);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitGroup(grp *GroupingExpr) (Value, error) {
	p.print_header("Group");
	p.tab();
		p.print_def_expr("InnerExpr", grp.InnerExpr);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitAssign(ass *AssignExpr) (Value, error) {
	p.print_header("Assign");
	p.tab();
		p.print_def_token("Name", ass.Name);
		p.print_def_expr("Name", ass.Asset);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitFuncCall(fnc *FuncCall) (Value, error) {
	p.print_header("FunctionCall");
	p.tab();
		p.print_def_expr("Calle", fnc.Callee);
		p.print_def_expr("Args", fnc.Args...);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitExpr(stmt *ExprStmt) (Value, error) {
	stmt.InnerExpr.Accept(p);
	return nil, nil;
}

func (p *PrettyPrinter) VisitVariableDeclaration(vard *VarDeclarationStmt) (Value, error) {
	p.print_header("VariableDeclaration");
	p.tab();
		p.print_def_token("Name", vard.Name);
		p.print_def_expr("Asset", vard.Asset);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitFuncDeclarationStmt(fnd *FuncDeclarationStmt) (Value, error) {
	p.print_header("FunctionDeclaration");
	p.tab();
		p.print_def_token("Name", fnd.Name);
		p.print_def_token("Params", fnd.Params...);
		p.print_def_stmt("Body", fnd.Body...);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitReturn(ret *ReturnStmt) (Value, error) {
	p.print_header("ReturnStamement");
	p.tab();
		p.print_def_expr("Asset", ret.Asset);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitPrint(prnt *PrintStmt) (Value, error) {
	p.print_header("PrintStatement");
	p.tab();
		p.print_def_expr("Asset", prnt.Asset);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitBlock(blk *BlockStmt) (Value, error) {
	p.print_header("BlockStatement");
	p.tab();
		p.print_def_stmt("Body", blk.Stmts...);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitBreak(*BreakStmt) (Value, error) {
	p.print_header("BreakStatement");
	return nil, nil;
}

func (p *PrettyPrinter) VisitContinue(*ContinueStmt) (Value, error) {
	p.print_header("ContinueStatement");
	return nil, nil;
}

func (p *PrettyPrinter) VisitConditional(cond *ConditionalStmt) (Value, error) {
	p.print_header("ConditionalStatment");
	p.tab();
		for _, b := range cond.Branches {
			p.print_def_expr("Cond", b.Condition);
			p.print_def_stmt("Body", b.NDStmt);
		}
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitWhile(whl *WhileStmt) (Value, error) {
	p.print_header("WhileStatement");
	p.tab();
		p.print_def_expr("Cond", whl.Cond);
		p.print_def_stmt("Body", whl.NDStmt);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) VisitFor(fors *ForStmt) (Value, error) {
	p.print_header("ForStatement");
	p.tab();
		if fors.Init != nil {
			p.print_def_stmt("Init", fors.Init);
		}
		if fors.Cond != nil {
			p.print_def_expr("Cond", fors.Cond);
		}
		if fors.Step != nil {
			p.print_def_expr("Step", fors.Step);
		}
		p.print_def_stmt("Body", fors.NDStmt);
	p.untab();
	return nil, nil;
}

func (p *PrettyPrinter) Print(stmt Stmt) {
	stmt.Accept(p);
	fmt.Print(p.builder.String());
}
