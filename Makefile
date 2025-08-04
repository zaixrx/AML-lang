all: file

file: main.out file.aml
	./main.out -mode file file.aml

repl: main.out
	./main.out -mode repl

main.out: *.go
	go build -o main.out && ./main.out
