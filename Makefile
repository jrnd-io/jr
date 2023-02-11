hello:
	@echo "JG, the JSON Generator"

compile:
	@echo "Compiling"
	go build -o build jg.go
	
compile-all:
	@echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o build/jg-linux-arm jg.go
	GOOS=linux GOARCH=arm64 go build -o build/jg-linux-arm64 jg.go
	GOOS=freebsd GOARCH=386 go build -o build/jg-freebsd-386 jg.go

run: compile
	./jg

clean:
	go clean
	rm build/*

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

help: hello
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}all${RESET}'
	@echo ''

all: hello compile
