package interpreter

import (
	"aml/parser"
	"fmt"
	"reflect"
);

type ReturnError struct {
	value parser.Value;
}

type BreakError struct {}

type ContinueError struct {}

func (err *ReturnError) Error() string {
	return "RUNTIME ERROR: 'return' should only be used inside a function";
}

func (err *BreakError) Error() string {
	return "RUNTIME ERROR: 'break' should only be used inside 'for' or 'while'";
}

func (err *ContinueError) Error() string {
	return "RUNTIME ERROR: 'continue' should only be used inside 'for' or 'while'";
}

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

func (in Interpreter) extract_number(value parser.Value) (bool, float64) {
	if number, ok := value.(float64); ok {
		return true, number;
	}
	return false, 0;
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

func (in Interpreter) extract_string(value parser.Value) string {
	return fmt.Sprint(value);
}

func (in Interpreter) equal(a parser.Value, b parser.Value) bool {
	return reflect.ValueOf(a).Equal(reflect.ValueOf(b));
}
