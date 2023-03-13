VERSION=0.0.2
USER=$(shell id -u -n)
TIME=$(shell date)

hello:
	@echo "JR,the JSON Random Generator"

compile:
	@echo "Compiling"
	go build -v -ldflags="-X 'jr/cmd.Version=$(VERSION)' -X 'jr/cmd.BuildUser=$(USER)' -X 'jr/cmd.BuildTime=$(TIME)'" -o build jr.go
	
compile-all:
	@echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o build/jg-linux-arm jr.go
	GOOS=linux GOARCH=arm64 go build -o build/jg-linux-arm64 jr.go
	GOOS=freebsd GOARCH=386 go build -o build/jg-freebsd-386 jr.go

run: compile
	./build/jr

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

copy_templates:
	mkdir -p ~/.jr && cp -r templates ~/.jr/

install_bin:
	install build/jr /usr/local/bin

install: install_bin copy_templates

all: hello compile
