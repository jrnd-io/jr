VERSION=0.3.8
USER=$(shell id -u -n)
TIME=$(shell date)

hello:
	@echo "JR,the JSON Random Generator"

generate:
	go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
	go generate pkg/generator/generate.go

compile:
	@echo "Compiling"
	go build -v -ldflags="-X 'github.com/ugol/jr/pkg/cmd.Version=$(VERSION)' -X 'github.com/ugol/jr/pkg/cmd.BuildUser=$(USER)' -X 'github.com/ugol/jr/pkg/cmd.BuildTime=$(TIME)'" -o build/jr jr.go

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
	mkdir -p ~/.jr/kafka && cp -r templates ~/.jr/ && cp -r pkg/producers/kafka/*.properties.example ~/.jr/kafka/

copy_config:
	mkdir -p ~/.jr && cp config/* ~/.jr/

install:
	install build/jr /usr/local/bin

all: hello generate compile


