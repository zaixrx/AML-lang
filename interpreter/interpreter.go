package interpreter

import (
	"fmt"
	"reflect"
	"aml/lexer"
	"aml/parser"
);

type Interpreter struct {
	environment map[string]parser.Value;
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
		fmt.Println("not a float", reflect.TypeOf(value), value);
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
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '+' must be numbers");
			}
			return nums[0] + nums[1], nil;
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
	name := string(expr.Name.Lexeme);
	value, exist := in.environment[name];
	if !exist {
		return nil, in.generate_error("variable %s is not declared", name);
	}
	return value, nil;
}

func (in Interpreter) VisitAssign(expr *parser.AssignExpr) (parser.Value, error) {
	if _, pres := in.environment[expr.To]; !pres { 
		return nil, in.generate_error("variable of name '%s' isn't declared!", expr.To);
	}
	value, err := expr.From.Accept(in);
	if err != nil {
		return nil, err;
	}
	in.environment[expr.To] = value;
	return value, nil;
}

// statements
func (in Interpreter) VisitExpr(stmt *parser.ExprStmt) error {
	_, err := stmt.InnerExpr.Accept(in);
	return err;
}

func (in Interpreter) VisitVariableDeclaration(stmt *parser.VarDeclarationStmt) error {
	_, exist := in.environment[stmt.Name];
	if exist {
		return in.generate_error("variable %s is already declared", stmt.Name);
	}
	value, err := stmt.Asset.Accept(in);
	if err != nil {
		return err;
	}
	in.environment[stmt.Name] = value;
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

func (in Interpreter) exec(stmt parser.Stmt) error {
	return stmt.Accept(in);
}

func NewInterpreter() Interpreter {
	return Interpreter {
		environment: make(map[string]parser.Value),
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
