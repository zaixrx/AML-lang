package interpreter

import (
	"aml/lexer"
	"aml/parser"
	"errors"
	"fmt"
	"strings"
)

type ReturnError struct {
	val parser.Value
}
func (e *ReturnError) Error() string {
	return "RUNTIME ERROR: 'return' should only be used inside a function";
}

var BreakError = fmt.Errorf("RUNTIME ERROR: 'break' should only be used inside 'for' or 'while'");
var ContinueError = fmt.Errorf("RUNTIME ERROR: 'continue' should only be used inside 'for' or 'while'");

type Environment struct {
	refs map[string]parser.Value;
	prev *Environment;
};

func NewEnvironment(prev *Environment) *Environment {
	return &Environment{
		refs: make(map[string]parser.Value),
		prev: prev,
	};
}

func (env *Environment) get(name string) (parser.Value, error) {
	curr := env;
	for curr != nil {
		if value, exists := curr.refs[name]; exists {
			return value, nil;
		}
		curr = curr.prev;
	}
	return nil, fmt.Errorf("variable %s is not declared", name);
}

func (env *Environment) declare(name string, value parser.Value) error {
	if _, exists := env.refs[name]; exists {
		return fmt.Errorf("variable %s is already declared", name);
	}
	env.refs[name] = value;
	return nil;
}

func (env *Environment) assign(name string, new_value parser.Value) error {
	curr := env;
	for curr != nil {
		if _, exists := curr.refs[name]; exists {
			curr.refs[name] = new_value;
			return nil;
		}
		curr = curr.prev;
	}
	return fmt.Errorf("variable %s is not declared", name);
}

type Callable interface {
	Arity() byte;
	Execute(in Interpreter, args []parser.Value) (parser.Value, error);
	String() string;
}

type Interpreter struct {
	environment *Environment;
};

func (in Interpreter) generate_error(format string, args ...any) error {
	return fmt.Errorf("RUNTIME ERROR: %s", fmt.Sprintf(format, args...));
}

func (in Interpreter) extract_boolean(value parser.Value) bool {
	if value == nil || value == false {
		return false;
	}
	return true;
}

func (in Interpreter) extract_string(value parser.Value) string {
	return fmt.Sprint(value);
}

func (in Interpreter) equal(a parser.Value, b parser.Value) bool {
	return a == b;
}

// expressions
func (in Interpreter) VisitUnary(expr parser.UnaryExpr) (parser.Value, error) {
	value, err := expr.Operand.Accept(in);
	if err != nil {
		return nil, err;
	}
	if value == nil {
		return nil, in.generate_error("cannot apply unary operator on null operand");
	}
	switch expr.Operator.Type {
		case lexer.BANG: {
			return !in.extract_boolean(value), nil;
		};
		case lexer.MINUS: {
			if num, ok := value.(float64); ok {
				return -num, nil;
			}
			return nil, in.generate_error("unary '-' can only be used on numbers");
		};
	}
	return nil, in.generate_error("invalid unary operation, got %s", expr.Operator.Type.ToString());
}

func (in Interpreter) VisitBinary(expr parser.BinaryExpr) (parser.Value, error) {
	leftval, err := expr.LOperand.Accept(in);
	if err != nil {
		return nil, err;
	}
	if leftval == nil {
		return nil, in.generate_error("cannot apply binary operator on null left operand");
	}
	rightval, err := expr.ROperand.Accept(in);
	if err != nil {
		return nil, err;
	}
	if rightval == nil {
		return nil, in.generate_error("cannot apply binary operator on null right operand");
	}
	switch expr.Operator.Type {
		case lexer.PLUS: {
			if lnum, ok := leftval.(float64); ok {
				if rnum, ok := rightval.(float64); ok {
					return lnum + rnum, nil;
				}
				return nil, in.generate_error("right operand in binary '+' must be number");
			} else if lstr, ok := leftval.(string); ok {
				if rstr, ok := rightval.(string); ok {
					return lstr + rstr, nil;
				}
				return nil, in.generate_error("right operand in binary '+' must be string");
			}
			return nil, in.generate_error("operands in binary '+' must be strings or numbers");
		};
		case lexer.MINUS: {
			if rnum, ok := rightval.(float64); ok {
				if lnum, ok := leftval.(float64); ok {
					return lnum - rnum, nil;
				}
				return nil, in.generate_error("right operand in binary '-' must be number");
			}
			return nil, in.generate_error("right operand in binary '*' must be number");
		};
		case lexer.STAR: {
			if rnum, ok := rightval.(float64); ok {
				if lnum, ok := leftval.(float64); ok {
					return lnum * rnum, nil;
				}
				return nil, in.generate_error("right operand in binary '*' must be number");
			}
			return nil, in.generate_error("right operand in binary '*' must be number");
		};
		case lexer.SLASH: {
			if rnum, ok := rightval.(float64); ok {
				if lnum, ok := leftval.(float64); ok {
					return lnum / rnum, nil;
				}
				return nil, in.generate_error("right operand in binary '/' must be number");
			}
			return nil, in.generate_error("right operand in binary '/' must be number");
		};
		case lexer.EQUAL_EQUAL: {
			return in.equal(leftval, rightval), nil;
		};
		case lexer.BANG_EQUAL: {
			return !in.equal(leftval, rightval), nil;
		};
		case lexer.GREATER: {
			if rnum, ok := rightval.(float64); ok {
				if lnum, ok := leftval.(float64); ok {
					return lnum > rnum, nil;
				}
				return nil, in.generate_error("right operand in binary '>' must be number");
			}
			return nil, in.generate_error("right operand in binary '>' must be number");
		};
		case lexer.GREATER_EQUAL: {
			if rnum, ok := rightval.(float64); ok {
				if lnum, ok := leftval.(float64); ok {
					return lnum >= rnum, nil;
				}
				return nil, in.generate_error("right operand in binary '>=' must be number");
			}
			return nil, in.generate_error("right operand in binary '>=' must be number");
		};
		case lexer.LESS: {
			if rnum, ok := rightval.(float64); ok {
				if lnum, ok := leftval.(float64); ok {
					return lnum < rnum, nil;
				}
				return nil, in.generate_error("right operand in binary '<' must be number");
			}
			return nil, in.generate_error("right operand in binary '<' must be number");
		};
		case lexer.LESS_EQUAL: {
			if rnum, ok := rightval.(float64); ok {
				if lnum, ok := leftval.(float64); ok {
					return lnum <= rnum, nil;
				}
				return nil, in.generate_error("right operand in binary '<=' must be number");
			}
			return nil, in.generate_error("right operand in binary '<=' must be number");
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

func (in Interpreter) VisitTernary(expr parser.TernaryExpr) (parser.Value, error) {
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

func (in Interpreter) VisitGroup(expr parser.GroupingExpr) (parser.Value, error) {
	return expr.InnerExpr.Accept(in);
}

func (in Interpreter) VisitLiteral(expr parser.LiteralExpr) (parser.Value, error) {
	if str, ok := expr.ValueLiteral.(string); ok {
		bytes := []byte(str);
		return string(bytes[1:len(bytes)-1]), nil;
	}
	return expr.ValueLiteral, nil;
}

func (in Interpreter) VisitVariable(expr parser.VariableExpr) (parser.Value, error) {
	value, err := in.environment.get(expr.Name.Lexeme);
	if err != nil {
		return nil, in.generate_error("%s", err.Error());
	}
	return value, nil;
}

func (in Interpreter) VisitAssign(expr parser.AssignExpr) (parser.Value, error) {
	value, err := expr.Asset.Accept(in);
	if err != nil {
		return nil, err;
	}
	err = in.environment.assign(expr.Name.Lexeme, value);
	if err != nil {
		return nil, err;
	}
	return value, nil;
}

func (in Interpreter) VisitFuncCall(expr parser.FuncCall) (parser.Value, error) {
	val, err := expr.Callee.Accept(in);
	if err != nil {
		return nil, err;
	}
	fn, callable_ok := val.(Callable);
	if !callable_ok {
		return nil, in.generate_error("invalid callee target");
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

func (in Interpreter) VisitReturn(stmt parser.ReturnStmt) (parser.Value, error) {
	var ( val parser.Value; err error = nil; );
	if stmt.Asset != nil {
		val, err = stmt.Asset.Accept(in);
		if err != nil {
			return nil, err;
		}
	}
	return nil, &ReturnError{ val: val };
}

func (in Interpreter) VisitBreak(_ parser.BreakStmt) (parser.Value, error) {
	return nil, BreakError;
}

func (in Interpreter) VisitContinue(_ parser.ContinueStmt) (parser.Value, error) {
	return nil, ContinueError;
}

// statements
func (in Interpreter) VisitExpr(stmt parser.ExprStmt) (parser.Value, error) {
	return stmt.InnerExpr.Accept(in);
}

func (in Interpreter) VisitVariableDeclaration(stmt parser.VarDeclarationStmt) (parser.Value, error) {
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
	err = in.environment.declare(stmt.Name.Lexeme, value);
	if err != nil {
		return nil, in.generate_error("%s", err.Error());
	}
	return nil, nil;
}

func (in Interpreter) VisitFuncDeclarationStmt(stmt parser.FuncDeclarationStmt) (parser.Value, error) {
	err := in.environment.declare(stmt.Name.Lexeme, AMLFunc{
		closure: in.environment,
		internal: parser.Func(stmt),
	});
	if err != nil {
		return nil, in.generate_error("%s", err.Error());
	}
	return nil, nil;
}

func (in Interpreter) VisitPrint(stmt parser.PrintStmt) (parser.Value, error) {
	builder := strings.Builder{};
	for i, asset := range stmt.Assets {
		val, err := asset.Accept(in);
		if err != nil {
			return nil, err;
		}
		if i != 0 {
			builder.WriteString(" ");
		}
		builder.WriteString(in.extract_string(val));
	}
	fmt.Println(builder.String());
	return nil, nil;
}

func (in Interpreter) VisitBlock(block parser.BlockStmt) (parser.Value, error) {
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

func (in Interpreter) VisitConditional(stmt parser.ConditionalStmt) (parser.Value, error) {
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

func (in Interpreter) VisitWhile(stmt parser.WhileStmt) (parser.Value, error) {
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
			if errors.Is(err, BreakError) {
				goto exit_loop;
			}
			if errors.Is(err, ContinueError) {
				return nil, err;
			}
		}
		goto loop;
	}
exit_loop:
	return val, nil;
}

func (in Interpreter) VisitFor(stmt parser.ForStmt) (parser.Value, error) {
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
			if errors.Is(err, BreakError) {
				goto exit_loop;
			}
			if errors.Is(err, ContinueError) {
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
