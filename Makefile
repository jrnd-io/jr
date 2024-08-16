VERSION=0.3.9
GOVERSION=$(shell go version)
USER=$(shell id -u -n)
TIME=$(shell date)

hello:
	@echo "JR,the JSON Random Generator"
    CONFIG_HOME=$XDG_CONFIG_HOME
    ifeq (CONFIG_HOME, "" )
      HOME = "~"
    endif
install-gogen:
	go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
	#go install github.com/hamba/avro/v2/cmd/avrogen@latest

generate:
	go generate pkg/generator/generate.go

compile:
	@echo "Compiling"
	go build -v -ldflags="-X 'github.com/ugol/jr/pkg/cmd.Version=$(VERSION)' -X 'github.com/ugol/jr/pkg/cmd.GoVersion=$(GOVERSION)' -X 'github.com/ugol/jr/pkg/cmd.BuildUser=$(USER)' -X 'github.com/ugol/jr/pkg/cmd.BuildTime=$(TIME)'" -o build/jr jr.go

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
	mkdir -p CONFIG_HOME/.jr/kafka && cp -r templates CONFIG_HOME/.jr/ && cp -r pkg/producers/kafka/*.properties.example CONFIG_HOME/.jr/kafka/

copy_config:
	mkdir -p CONFIG_HOME/.jr && cp config/* CONFIG_HOME/.jr/

install:
	install build/jr /usr/local/bin

all: hello install-gogen generate compile
all_offline: hello generate compile