package main

import (
	"bufio"
	"fmt"
	"os"
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
		for _, ast := range asts {
			fmt.Println(ast.String());
		}
	}
}
