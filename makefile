.SILENT: run
run:
	tput setaf 2; echo "tests:"; tput sgr0
	go test ./src/lexer
	go test ./src/ast
	go test ./src/parser
	go test ./src/object
	go test ./src/evaluator
	go test ./src/code
	echo
	tput setaf 2; echo "exec:"; tput sgr0
	go run src/main.go
