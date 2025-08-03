package main

import (
	"fmt"
	"reflect"
)

type Value = any

type Interpreter struct { };

func (in Interpreter) generate_error(format string, args ...any) error {
	return fmt.Errorf("RUNTIME ERROR: %s", fmt.Sprintf(format, args...));
}

func (in Interpreter) extract_boolean(value Value) bool {
	if value == nil {
		return false;
	}
	return true;
}

func (in Interpreter) extract_number(value Value) (bool, float64) {
	if reflect.TypeOf(value).Kind() != reflect.Float64 {
		fmt.Println("not a float", reflect.TypeOf(value), value);
		return false, 0;
	}
	return true, reflect.ValueOf(value).Float();
}

func (in Interpreter) extract_numbers(values ...Value) (bool, []float64) {
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

func (in Interpreter) equal(a Value, b Value) bool {
	return reflect.ValueOf(a).Equal(reflect.ValueOf(b));
}

func (in Interpreter) visit_unary(expr *UnaryExpr) (Value, error) {
	value, err := expr.operand.AcceptInterpreter(in);
	if err != nil {
		return nil, err;
	}
	switch expr.operator.Type {
		case BANG: {
			return !in.extract_boolean(value), nil;
		};
		case MINUS: {
			succ, num := in.extract_number(value);
			if !succ {
				return nil, in.generate_error("unary '-' can only be used on numbers");
			}
			return -num, nil;
		};
	}
	return nil, in.generate_error("invalid unary operation, got %s", expr.operator.Type.ToString());
}

func (in Interpreter) visit_binary(expr *BinaryExpr) (Value, error) {
	leftval, err := expr.left.AcceptInterpreter(in);
	if err != nil {
		return nil, err;
	}
	rightval, err := expr.right.AcceptInterpreter(in);
	if err != nil {
		return nil, err;
	}
	switch expr.operator.Type {
		case PLUS: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '+' must be numbers");
			}
			return nums[0] + nums[1], nil;
		};
		case MINUS: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '-' must be numbers");
			}
			return nums[0] - nums[1], nil;
		};
		case STAR: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '*' must be numbers");
			}
			return nums[0] * nums[1], nil;
		};
		case SLASH: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '/' must be numbers");
			}
			return nums[0] / nums[1], nil;
		};
		case EQUAL_EQUAL: {
			return in.equal(leftval, rightval), nil;
		};
		case BANG_EQUAL: {
			return !in.equal(leftval, rightval), nil;
		};
		case GREATER: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '>' must be numbers");
			}
			return nums[0] > nums[1], nil;
		};
		case GREATER_EQUAL: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '>=' must be numbers");
			}
			return nums[0] >= nums[1], nil;
		};
		case LESS: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '<' must be numbers");
			}
			return nums[0] < nums[1], nil;
		};
		case LESS_EQUAL: {
			succ, nums := in.extract_numbers(leftval, rightval);
			if !succ {
				return nil, in.generate_error("operands in binary '<=' must be numbers");
			}
			return nums[0] <= nums[1], nil;
		};
		case AND: {
			return in.extract_boolean(leftval) && in.extract_boolean(rightval), nil;
		};
		case OR: {
			return in.extract_boolean(leftval) || in.extract_boolean(rightval), nil;
		};
	}
	return nil, in.generate_error("invalid binary operation, got %s", expr.operator.Type.ToString());
}

func (in Interpreter) visit_ternary(expr *TernaryExpr) (Value, error) {
	condval, err := expr.cond.AcceptInterpreter(in);
	if err != nil {
		return nil, err;
	}
	if in.extract_boolean(condval) {
		value, err := expr.iftrue.AcceptInterpreter(in);
		if err != nil {
			return nil, err;
		}
		return value, nil;
	}
	value, err := expr.iffalse.AcceptInterpreter(in);
	if err != nil {
		return nil, err;
	}
	return value, nil;
}

func (in Interpreter) visit_group(expr *GroupingExpr) (Value, error) {
	return expr.expr.AcceptInterpreter(in);
}

func (in Interpreter) visit_literal(expr *LiteralExpr) (Value, error) {
	return expr.value, nil;
}

func (in Interpreter) Interpret(exprs []Expr) ([]Value, error) {
	values := make([]Value, len(exprs));
	for i, expr := range exprs {
		value, err := expr.AcceptInterpreter(in);
		if err != nil {
			return nil, err;
		}
		values[i] = value;
	}
	return values, nil;
}
