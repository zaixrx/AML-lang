package main

import (
	"os"
	"fmt"
	"flag"
	"bufio"

	"aml/lexer"
	"aml/parser"
	"aml/interpreter"
)

func evalAML(filename string, content string) {
	scanner := lexer.NewScanner(filename, content);
	tokens, err := scanner.Scan();
	if err != nil {
		fmt.Println(err);
		return;
	}
	
	for _, token := range tokens {
		fmt.Println(token);
	}
	
	parser := parser.NewParser(tokens);
	stmts, err := parser.Parse();
	if err != nil {
		fmt.Println(err);
		return;
	}

	for _, stmt := range stmts {
		fmt.Println(stmt);
	}

	interpreter := interpreter.Interpreter{};
	err = interpreter.Interpret(stmts);
	if err != nil {
		fmt.Println(err);
		return;
	}
}

func handleREPL() {
	reader := bufio.NewReader(os.Stdin);
	for {
		fmt.Print(">> ");
		code, err := reader.ReadString('\n');
		if err != nil {
			fmt.Println(err);
			fmt.Println("Terminating REPL Process...");
			break;
		}
		evalAML("REPL", code);
	}
}

func handleFile(filename string) error {
	bcode, err := os.ReadFile(filename);
	if err != nil {
		return err;
	}
	evalAML(filename, string(bcode));
	return nil;
}

func main() {
	mode := flag.String("mode", "repl", "set interpreter mode");
	flag.Parse();
	switch *mode {
		case "repl": {
			handleREPL();
		}
		case "file": {
			if len(os.Args) != 4 {
				fmt.Printf("usage: %s -mode file <file_name>\n", os.Args[1]);
				return;
			}
			err := handleFile(os.Args[3]);
			if err != nil {
				fmt.Println(err);
			}
		}
		default: {
			fmt.Println("unsupported mode");
		}
	};
}
