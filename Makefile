VERSION=0.3.9
GOVERSION=$(shell go version)
USER=$(shell id -u -n)
TIME=$(shell date)
JR_HOME=jr

ifndef XDG_CONFIG_HOME
ifeq ($(OS), Windows_NT)
	detectedOS := Windows
else
	detectedOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

ifeq ($(detectedOS), Darwin)
	CONFIG_HOME="$(HOME)/Library/Application Support"
endif
ifeq ($(detectedOS),  Linux)
	CONFIG_HOME="$(HOME)/.config"
endif
ifeq ($(detectedOS), Windows_NT)
	CONFIG_HOME="$(LOCALAPPDATA)"
endif
else
	CONFIG_HOME=$(XDG_CONFIG_HOME)
endif

hello:
	@echo "JR,the JSON Random Generator"
	@echo " Version: $(VERSION)"
	@echo " Go Version: $(GOVERSION)"
	@echo " Build User: $(USER)"
	@echo " Build Time: $(TIME)"
	@echo " Detected OS: $(detectedOS)"
	@echo " Config Home: $(CONFIG_HOME)"

install-gogen:
	go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
	#go install github.com/hamba/avro/v2/cmd/avrogen@latest

generate:
	go generate pkg/generator/generate.go

compile:
	@echo "Compiling"
	go build -v -ldflags="-X 'github.com/ugol/jr/pkg/cmd.Version=$(VERSION)' \
	-X 'github.com/ugol/jr/pkg/cmd.GoVersion=$(GOVERSION)' \
	-X 'github.com/ugol/jr/pkg/cmd.BuildUser=$(USER)' \
	-X 'github.com/ugol/jr/pkg/cmd.BuildTime=$(TIME)'" \
	-o build/jr jr.go

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
	mkdir -p $(CONFIG_HOME)/$(JR_HOME)/kafka && \
	cp -r templates $(CONFIG_HOME)/$(JR_HOME) && \
	cp -r pkg/producers/kafka/*.properties.example $(CONFIG_HOME)/$(JR_HOME)/kafka/

copy_config:
	mkdir -p $(CONFIG_HOME)/$(JR_HOME) && \
	cp config/* $(CONFIG_HOME)/$(JR_HOME)/

install:
	install build/jr /usr/local/bin

all: hello install-gogen generate compile
all_offline: hello generate compile
