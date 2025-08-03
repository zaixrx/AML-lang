package main

import (
	"os"
	"fmt"
	"bufio"
)

func main() {
	reader:= bufio.NewReader(os.Stdin);
	for {
		fmt.Print(">> ");
		str, err := reader.ReadString('\n');
		if err != nil {
			fmt.Println(err);
			fmt.Println("Terminating REPL Process...");
			break;
		}
		scanner := NewScanner("REPL", str);
		tokens, err := scanner.Scan();
		if err != nil {
			fmt.Println(err);
			continue;
		}
		parser := NewParser(tokens);
		asts, err := parser.Parse();
		if err != nil {
			fmt.Println(err);
			continue;
		}
		interpreter := Interpreter{};
		values, err := interpreter.Interpret(asts);
		if err != nil {
			fmt.Println(err);
			continue;
		}
		for _, value := range values {
			fmt.Printf("%#v\n", value);
		}
	}
}
