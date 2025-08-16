package interpreter

import (
	"aml/parser"
	"bufio"
	"os"
	"time"
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
	return string(bytes[0:len(bytes)-1]), nil;
}

func (StdRead) String() string {
	return "native: stdread/0";
}

type StdTime struct {};

func (StdTime) Arity() byte {
	return 0;
}

func (StdTime) Execute(Interpreter, []parser.Value) (parser.Value, error) {
	return float64(time.Now().UnixMilli()), nil;
}

func (StdTime) String() string {
	return "native: stdtime/0";
}

func GetStdFuncs() map[string]parser.Value {
	return map[string]parser.Value{
		"read": StdRead{},
		"time": StdTime{},
	};
}
