package interpreter

import (
	"fmt"
	"errors"

	"aml/parser"
);

type AMLFunc struct {
	closure *Environment;
	internal parser.Func;
}

func (fn AMLFunc) Arity() byte {
	return byte(len(fn.internal.Params));
}

func (fn AMLFunc) Execute(in Interpreter, args []parser.Value) (parser.Value, error) {
	old_env := in.environment;
	env := NewEnvironment(fn.closure); in.environment = env;
	defer func() { in.environment = old_env; }();
	for i, arg := range args {
		err := env.declare(fn.internal.Params[i].Lexeme, arg);
		if err != nil {
			return nil, err;
		}
	}
	var (
		reterr *ReturnError = nil
		retvalue parser.Value = nil
	);
	for _, stmt := range fn.internal.Body {
		if _, err := stmt.Accept(in); err != nil {
			if errors.As(err, &reterr) {
				retvalue = reterr.val;
				break;
			}
			return nil, err;
		}
	}
	return retvalue, nil;
}

func (fn AMLFunc) String() string {
	return fmt.Sprintf("function %s/%d", fn.internal.Name.Lexeme, fn.Arity());
}
