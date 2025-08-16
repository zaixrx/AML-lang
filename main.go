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

func evalAML(interpreter *interpreter.Interpreter, filename string, content string, use_pp bool) parser.Value {
	s := lexer.NewScanner(filename, content);
	tokens, err := s.Scan();
	if err != nil {
		fmt.Println(err);
		return nil;
	}
	// for _, token := range tokens {
	// 	fmt.Println(token);
	// }
	p := parser.NewParser(tokens);
	stmts, err := p.Parse();
	if err != nil {
		fmt.Println(err);
		return nil;
	}
	if use_pp {
		pp := parser.PrettyPrinter{};
		for _, stmt := range stmts {
		 	pp.Print(stmt);
		}
	}
	val, err := interpreter.Interpret(stmts);
	if err != nil {
		fmt.Println(err);
		return nil;
	}
	return val;
}

func handleREPL(use_pp bool) {
	reader := bufio.NewReader(os.Stdin);
	i := interpreter.NewInterpreter(); 
	for {
		fmt.Print(">> ");
		code, err := reader.ReadString('\n');
		if err != nil {
			fmt.Println(err);
			fmt.Println("Terminating REPL Process...");
			break;
		}
		val := evalAML(&i, "REPL", code, use_pp);
		if val != nil {
			fmt.Println(val);
		}
	}
}

func handleFile(filename string, use_pp bool) error {
	bcode, err := os.ReadFile(filename);
	if err != nil {
		return err;
	}
	i := interpreter.NewInterpreter();
	evalAML(&i, filename, string(bcode), use_pp);
	return nil;
}

func main() {
	repl := flag.Bool("repl", false, "use repl? else interpret file")
	use_pp := flag.Bool("p", true, "Use PrettyPrinter to Print ASTs");
	flag.Parse();
	if *repl {
		handleREPL(*use_pp);
	} else {
		if len(os.Args) < 2 {
			fmt.Printf("usage: %s <file_name> [-p]\n", os.Args[0]);
			return;
		}
		err := handleFile(os.Args[1], *use_pp);
		if err != nil {
			fmt.Println(err);
		}
	}
}
