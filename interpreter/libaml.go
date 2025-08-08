package interpreter

import (
	"os"
	"bufio"
	"aml/parser"
);

type StdRead struct {};

func (StdRead) Arity() byte {
	return 0;
}

func (StdRead) Execute(Interpreter, []parser.Value) (parser.Value, error) {
	reader := bufio.NewReader(os.Stdin);
	bytes, err := reader.ReadString('\n');
	if err != nil {
		return "", err;
	}
	return bytes, nil;
}

func GetStdFuncs() map[string]parser.Value {
	return map[string]parser.Value{
		"read": StdRead{},
	};
}
