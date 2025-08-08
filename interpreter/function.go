package interpreter

import (
	"aml/parser"
	"errors"
);

type AMLFunc parser.Func;

func (fn AMLFunc) Arity() byte {
	return byte(len(fn.Params));
}

func (fn AMLFunc) Execute(in Interpreter, args []parser.Value) (parser.Value, error) {
	env := NewEnvironment(in.environment);
	in.environment = env;
	defer func() { in.environment = in.environment.prev }();
	for i, arg := range args {
		err := env.declare(string(fn.Params[i].Lexeme), arg);
		if err != nil {
			return nil, err;
		}
	}
	var (
		retvalue parser.Value = nil
		reterr *ReturnError = nil
	);
	for _, stmt := range fn.Body {
		if _, err := stmt.Accept(in); err != nil {
			if errors.As(err, &reterr) {
				retvalue = reterr.value;
				break;
			}
			return nil, err;
		}
	}
	return retvalue, nil;
}
