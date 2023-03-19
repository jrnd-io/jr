VERSION=0.0.7
USER=$(shell id -u -n)
TIME=$(shell date)

hello:
	@echo "JR,the JSON Random Generator"

compile:
	@echo "Compiling"
	go build -v -ldflags="-X 'jr/cmd.Version=$(VERSION)' -X 'jr/cmd.BuildUser=$(USER)' -X 'jr/cmd.BuildTime=$(TIME)'" -o build/jr jr.go
	

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

install:
	install build/jr /usr/local/bin

all: hello compile
