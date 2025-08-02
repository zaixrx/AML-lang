package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if (len(os.Args) != 2) {
		fmt.Printf("usage: %s <file.aml>\n", os.Args[0]);
		return;
	}
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
		for _, token := range tokens {
			fmt.Printf("Token: {\n    lexme: \"%s\",\n   literal: %v\n    type: %s\n}\n", string(token.Lexeme), token.Literal, token.Type.ToString());
		}
	}
}
