test:
	go test ./...

trace:
	go test -v -run TestOperatorPrecedenceParsing ./parser

.PHONY: test trace