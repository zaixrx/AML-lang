package interpreter

import (
	"fmt"
	"errors"
	"reflect"
	"aml/lexer"
	"aml/parser"
);

// expressions
func (in Interpreter) VisitUnary(expr *parser.UnaryExpr) (parser.Value, error) {
	value, err := expr.Operand.Accept(in);
	if err != nil {
		return nil, err;
	}
	switch expr.Operator.Type {
		case lexer.BANG: {
			return !in.extract_boolean(value), nil;
		};
		case lexer.MINUS: {
			succ, num := in.extract_number(value);
			if !succ {
				return nil, in.generate_error("unary '-' can only be used on numbers");
			}
			return -num, nil;
		};
	}
	return nil, in.generate_error("invalid unary operation, got %s", expr.Operator.Type.ToString());
}

func (in Interpreter) VisitBinary(expr *parser.BinaryExpr) (parser.Value, error) {
	leftval, err := expr.LOperand.Accept(in);
	if err != nil {
		return nil, err;
	}
	rightval, err := expr.ROperand.Accept(in);
	if err != nil {
		return nil, err;
	}
	switch expr.Operator.Type {
		case lexer.PLUS: {
			if reflect.TypeOf(leftval).Kind() == reflect.Float64 {
				succ, nums := in.extract_numbers(leftval, rightval);
				if !succ {
					return nil, in.generate_error("operands in binary '+' must be numbers");
				}
				return nums[0] + nums[1], nil;
			} else if reflect.TypeOf(leftval).Kind() == reflect.String {
				return string(in.extract_string(leftval) + in.extract_string(rightval)), nil;
			}
			return nil, in.generate_error("operands in binary '+' must be strings or numbers");
		};
		case lexer.MINUS: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '-' must be numbers");
			}
			return nums[0] - nums[1], nil;
		};
		case lexer.STAR: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '*' must be numbers");
			}
			return nums[0] * nums[1], nil;
		};
		case lexer.SLASH: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '/' must be numbers");
			}
			return nums[0] / nums[1], nil;
		};
		case lexer.EQUAL_EQUAL: {
			return in.equal(leftval, rightval), nil;
		};
		case lexer.BANG_EQUAL: {
			return !in.equal(leftval, rightval), nil;
		};
		case lexer.GREATER: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '>' must be numbers");
			}
			return nums[0] > nums[1], nil;
		};
		case lexer.GREATER_EQUAL: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '>=' must be numbers");
			}
			return nums[0] >= nums[1], nil;
		};
		case lexer.LESS: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '<' must be numbers");
			}
			return nums[0] < nums[1], nil;
		};
		case lexer.LESS_EQUAL: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '<=' must be numbers");
			}
			return nums[0] <= nums[1], nil;
		};
		case lexer.AND: {
			return in.extract_boolean(leftval) && in.extract_boolean(rightval), nil;
		};
		case lexer.OR: {
			return in.extract_boolean(leftval) || in.extract_boolean(rightval), nil;
		};
	}
	return nil, in.generate_error("invalid binary operation, got %s", expr.Operator.Type.ToString());
}

func (in Interpreter) VisitTernary(expr *parser.TernaryExpr) (parser.Value, error) {
	condval, err := expr.Cond.Accept(in);
	if err != nil {
		return nil, err;
	}
	if in.extract_boolean(condval) {
		value, err := expr.Iftrue.Accept(in);
		if err != nil {
			return nil, err;
		}
		return value, nil;
	}
	value, err := expr.Iffalse.Accept(in);
	if err != nil {
		return nil, err;
	}
	return value, nil;
}

func (in Interpreter) VisitGroup(expr *parser.GroupingExpr) (parser.Value, error) {
	return expr.InnerExpr.Accept(in);
}

func (in Interpreter) VisitLiteral(expr *parser.LiteralExpr) (parser.Value, error) {
	return expr.ValueLiteral, nil;
}

func (in Interpreter) VisitVariable(expr *parser.VariableExpr) (parser.Value, error) {
	value, err := in.environment.get(string(expr.Name.Lexeme));
	if err != nil {
		return nil, in.generate_error("%s", err.Error());
	}
	return value, nil;
}

func (in Interpreter) VisitAssign(expr *parser.AssignExpr) (parser.Value, error) {
	value, err := expr.Asset.Accept(in);
	if err != nil {
		return nil, err;
	}
	err = in.environment.assign(expr.Name, value);
	if err != nil {
		return nil, err;
	}
	return value, nil;
}

func (in Interpreter) VisitFuncCall(expr *parser.FuncCall) (parser.Value, error) {
	val, err := expr.Callee.Accept(in);
	if err != nil {
		return nil, err;
	}
	fn, callable_ok := val.(Callable); // ok is only true for foreign functions
	if !callable_ok {
		parser_fn, func_ok := val.(parser.Func);
		if !func_ok {
			return nil, in.generate_error("invalid callee target");
		}
		fn = AMLFunc(parser_fn); // UGLY
	}
	if int(fn.Arity()) != len(expr.Args) {
		return nil, in.generate_error("expected %d arguments got %d", fn.Arity(), len(expr.Args));
	}
	args := make([]parser.Value, fn.Arity());
	for i := range fn.Arity() {
		val, err := expr.Args[i].Accept(in);
		if err != nil {
			return nil, err;
		}
		args[i] = val;
	}
	return fn.Execute(in, args);
}

func (in Interpreter) VisitReturn(stmt *parser.ReturnStmt) (parser.Value, error) {
	var ( value parser.Value; err error = nil; );
	if stmt.Asset != nil {
		value, err = stmt.Asset.Accept(in);
		if err != nil {
			return nil, err;
		}
	}
	return nil, &ReturnError{
		value: value,
	};
}

func (in Interpreter) VisitBreak(_ *parser.BreakStmt) (parser.Value, error) {
	return nil, &BreakError{};
}

func (in Interpreter) VisitContinue(_ *parser.ContinueStmt) (parser.Value, error) {
	return nil, &ContinueError{};
}

// statements
func (in Interpreter) VisitExpr(stmt *parser.ExprStmt) (parser.Value, error) {
	return stmt.InnerExpr.Accept(in);
}

func (in Interpreter) VisitVariableDeclaration(stmt *parser.VarDeclarationStmt) (parser.Value, error) {
	var (
		err error
		value parser.Value = nil;
	);
	if stmt.Asset != nil {
		value, err = stmt.Asset.Accept(in);
		if err != nil {
			return nil, err;
		}
	}
	err = in.environment.declare(stmt.Name, value);
	if err != nil {
		return nil, in.generate_error("%s", err.Error());
	}
	return nil, nil;
}

func (in Interpreter) VisitFuncDeclarationStmt(stmt *parser.FuncDeclarationStmt) (parser.Value, error) {
	err := in.environment.declare(string(stmt.Name.Lexeme), parser.Func(*stmt));
	if err != nil {
		return nil, in.generate_error("%s", err.Error());
	}
	return nil, nil;
}

func (in Interpreter) VisitPrint(stmt *parser.PrintStmt) (parser.Value, error) {
	val, err := stmt.Asset.Accept(in);
	if err != nil {
		return nil, err;
	}
	fmt.Printf("%s\n", in.extract_string(val));
	return nil, nil;
}

func (in Interpreter) VisitBlock(block *parser.BlockStmt) (parser.Value, error) {
	var (
		val parser.Value = nil;
		env = NewEnvironment(in.environment);
	);
	// create new environment
	in.environment = env;
	defer func () {
		in.environment = in.environment.prev;
	}();
	// execute all environment statements
	for _, stmt := range block.Stmts {
		sval, err := stmt.Accept(in);
		if err != nil {
			return nil, err;
		}
		if sval != nil {
			val = sval;
		}
	}
	return val, nil;
}

func (in Interpreter) VisitConditional(stmt *parser.ConditionalStmt) (parser.Value, error) {
	for _, branch := range stmt.Branches {
		if branch.Condition != nil {
			val, err := branch.Condition.Accept(in);
			if err != nil {
				return nil, err;
			}
			if !in.extract_boolean(val) {
				continue;
			}
		}
		_, err := branch.NDStmt.Accept(in);
		if err != nil {
			return nil, err;
		}
		break;
	}
	return nil, nil;
}

func (in Interpreter) VisitWhile(stmt *parser.WhileStmt) (parser.Value, error) {
	var (
		err error = nil;
		val parser.Value = nil;
	);
loop:
	cond, err := stmt.Cond.Accept(in)
	if err != nil {
		return nil, err;
	}
	if in.extract_boolean(cond) {
		val, err = stmt.NDStmt.Accept(in);
		if err != nil {
			return nil, err;
		}
		goto loop;
	}
	return val, nil;
}

func (in Interpreter) VisitFor(stmt *parser.ForStmt) (parser.Value, error) {
	var (
		err error = nil;
		val parser.Value = nil;
		cond parser.Value = true;
	);
	if stmt.Init != nil {
		_, err = stmt.Init.Accept(in);
		if err != nil {
			return nil, err;
		}
	}
	var (
		break_err *BreakError;
		continue_err *ContinueError;
	);
loop:
	if stmt.Cond != nil {
		cond, err = stmt.Cond.Accept(in)
		if err != nil {
			return nil, err;
		}
	}
	if in.extract_boolean(cond) {
		val, err = stmt.NDStmt.Accept(in);
		if err != nil {
			if errors.As(err, &break_err) {
				goto exit_loop;
			}
			// continue will basically continue
			if !errors.As(err, &continue_err) {
				return nil, err;
			}
		}
		if stmt.Step != nil {
			_, err = stmt.Step.Accept(in);
			if err != nil {
				return nil, err;
			}
		}
		goto loop;
	}
exit_loop:
	return val, nil;
}

func NewInterpreter() Interpreter {
	global_env := NewEnvironment(nil);
	for key, val := range GetStdFuncs() {
		global_env.declare(key, val);
	}
	return Interpreter {
		environment: global_env,
	};
}

func (in Interpreter) Interpret(stmts []parser.Stmt) (parser.Value, error) {
	var val parser.Value;
	for _, stmt := range stmts {
		sval, err := stmt.Accept(in);
		if err != nil {
			return nil, err;
		}
		if sval != nil {
			val = sval;
		}
	}
	return val, nil;
}
