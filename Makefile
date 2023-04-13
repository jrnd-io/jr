VERSION=0.1.6
USER=$(shell id -u -n)
TIME=$(shell date)

hello:
	@echo "JR,the JSON Random Generator"

generate:
	go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
	go install cmd/register/register.go
	go generate types/generate.go

compile:
	@echo "Compiling"
	go build -v -ldflags="-X 'github.com/ugol/jr/cmd.Version=$(VERSION)' -X 'github.com/ugol/jr/cmd.BuildUser=$(USER)' -X 'github.com/ugol/jr/cmd.BuildTime=$(TIME)'" -o build/jr jr.go

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
	mkdir -p ~/.jr/kafka && cp -r templates ~/.jr/ && cp -r producers/kafka/*.properties.example ~/.jr/kafka/

install:
	install build/jr /usr/local/bin

all: hello generate compile


