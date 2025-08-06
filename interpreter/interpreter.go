package interpreter

import (
	"fmt"
	"reflect"
	"aml/lexer"
	"aml/parser"
);

type Environment struct {
	vars map[string]parser.Value;
	prev *Environment;
};

func NewEnvironment() *Environment {
	return &Environment{
		vars: make(map[string]parser.Value),
		prev: nil,
	};
}

func (env *Environment) get(name string) (parser.Value, error) {
	curr := env;
	for curr != nil {
		if value, exists := curr.vars[name]; exists {
			return value, nil;
		}
		curr = curr.prev;
	}
	return nil, fmt.Errorf("variable %s is not declared", name);
}

func (env *Environment) declare(name string, initial_value parser.Value) error {
	if _, exists := env.vars[name]; exists {
		return fmt.Errorf("variable %s is already declared", name);
	}
	env.vars[name] = initial_value;
	return nil;
}

func (env *Environment) assign(name string, new_value parser.Value) error {
	curr := env;
	for curr != nil {
		if _, exists := curr.vars[name]; exists {
			curr.vars[name] = new_value;
			return nil;
		}
		curr = curr.prev;
	}
	return fmt.Errorf("variable %s is not declared", name);
}

type Interpreter struct {
	environment *Environment;
};

func (in Interpreter) generate_error(format string, args ...any) error {
	return fmt.Errorf("RUNTIME ERROR: %s", fmt.Sprintf(format, args...));
}

func (in Interpreter) extract_boolean(value parser.Value) bool {
	if value == nil {
		return false;
	}
	return true;
}

func (in Interpreter) extract_number(value parser.Value) (bool, float64) {
	if reflect.TypeOf(value).Kind() != reflect.Float64 {
		return false, 0;
	}
	return true, reflect.ValueOf(value).Float();
}

func (in Interpreter) extract_numbers(values ...parser.Value) (bool, []float64) {
	nums := make([]float64, len(values));
	for i, value := range values {
		succ, num := in.extract_number(value);
		if !succ {
			return false, nil;
		}
		nums[i] = num;
	}
	return true, nums;
}

func (in Interpreter) equal(a parser.Value, b parser.Value) bool {
	return reflect.ValueOf(a).Equal(reflect.ValueOf(b));
}

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

func (in Interpreter) stringify(value parser.Value) string {
	return fmt.Sprint(value);
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
				return string(in.stringify(leftval) + in.stringify(rightval)), nil;
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
		return nil, in.generate_error("%s", err.Error());
	}
	return value, nil;
}

// statements
func (in Interpreter) VisitExpr(stmt *parser.ExprStmt) error {
	_, err := stmt.InnerExpr.Accept(in);
	return err;
}

func (in Interpreter) VisitVariableDeclaration(stmt *parser.VarDeclarationStmt) error {
	initial_value, err := stmt.Asset.Accept(in);
	if err != nil {
		return err;
	}
	err = in.environment.declare(stmt.Name, initial_value);
	if err != nil {
		return in.generate_error("%s", err.Error());
	}
	return nil;
}

func (in Interpreter) VisitPrint(stmt *parser.PrintStmt) error {
	val, err := stmt.Asset.Accept(in);
	if err != nil {
		return err;
	}
	fmt.Printf("%v\n", val);
	return nil;
}

func (in Interpreter) VisitBlock(block *parser.BlockStmt) error {
	var err error = nil;
	// create new environment
	env := NewEnvironment();
	env.prev = in.environment;
	in.environment = env;
	// execute all environment statements
	for _, stmt := range block.Stmts {
		err = stmt.Accept(in);
		if err != nil {
			break;
		}
	}
	// return to old environment
	in.environment = in.environment.prev;
	return err;
}

func (in Interpreter) exec(stmt parser.Stmt) error {
	return stmt.Accept(in);
}

func NewInterpreter() Interpreter {
	return Interpreter {
		environment: NewEnvironment(),
	};
}

func (in Interpreter) Interpret(stmts []parser.Stmt) error {
	for _, stmt := range stmts {
		err := in.exec(stmt);
		if err != nil {
			return err;
		}
	}
	return nil;
}
