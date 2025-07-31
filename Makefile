DEFAULT_CONTENT=// this is a comment \
		(( )){} // grouping stuff \
		!*+-/=<> <= == // operators \

all: main.out file.amel
	./main.out file.amel

main.out: main.go
	go build -o main.out

file.amel:
	echo "$(DEFAULT_CONTENT)" > file.amel
