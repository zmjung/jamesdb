.DEFAULT_GOAL := refresh

fmt:
	@gofmt -s -w .
	@echo Formatted code successfully.

tidy:
	@go mod tidy
	@echo Tidied go.mod and go.sum files.

refresh: fmt tidy
