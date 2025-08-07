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

func evalAML(interpreter *interpreter.Interpreter, filename string, content string) parser.Value {
	scanner := lexer.NewScanner(filename, content);
	tokens, err := scanner.Scan();
	if err != nil {
		fmt.Println(err);
		return nil;
	}
	// for _, token := range tokens {
	// 	fmt.Println(token);
	// }
	parser := parser.NewParser(tokens);
	stmts, err := parser.Parse();
	if err != nil {
		fmt.Println(err);
		return nil;
	}
	val, err := interpreter.Interpret(stmts);
	if err != nil {
		fmt.Println(err);
		return nil;
	}
	return val;
}

func handleREPL() {
	reader := bufio.NewReader(os.Stdin);
	interpreter := interpreter.NewInterpreter(); 
	for {
		fmt.Print(">> ");
		code, err := reader.ReadString('\n');
		if err != nil {
			fmt.Println(err);
			fmt.Println("Terminating REPL Process...");
			break;
		}
		val := evalAML(&interpreter, "REPL", code);
		if val != nil {
			fmt.Println(val);
		}
	}
}

func handleFile(filename string) error {
	bcode, err := os.ReadFile(filename);
	if err != nil {
		return err;
	}
	interpreter := interpreter.NewInterpreter();
	evalAML(&interpreter, filename, string(bcode));
	return nil;
}

func main() {
	repl := flag.Bool("repl", false, "use repl? else interpret file")
	flag.Parse();
	if *repl {
		handleREPL();
	} else {
		if len(os.Args) != 2 {
			fmt.Printf("usage: %s <file_name>\n", os.Args[0]);
			return;
		}
		err := handleFile(os.Args[1]);
		if err != nil {
			fmt.Println(err);
		}
	}
}
