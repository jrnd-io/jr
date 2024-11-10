VERSION=0.4.0
GOVERSION=$(shell go version)
USER=$(shell id -u -n)
TIME=$(shell date)
JR_HOME=jr

ifndef XDG_DATA_DIRS
ifeq ($(OS), Windows_NT)
	detectedOS := Windows
else
	detectedOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

ifeq ($(detectedOS), Darwin)
	JR_SYSTEM_DIR="$(HOME)/Library/Application Support"
endif
ifeq ($(detectedOS),  Linux)
	JR_SYSTEM_DIR="$(HOME)/.config"
endif
ifeq ($(detectedOS), Windows_NT)
	JR_SYSTEM_DIR="$(LOCALAPPDATA)"
endif
else
	JR_SYSTEM_DIR=$(XDG_DATA_DIRS)
endif

ifndef XDG_DATA_HOME
ifeq ($(OS), Windows_NT)
	detectedOS := Windows
else
	detectedOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

ifeq ($(detectedOS), Darwin)
	JR_USER_DIR="$(HOME)/.local/share"
endif
ifeq ($(detectedOS),  Linux)
	JR_USER_DIR="$(HOME)/.local/share"
endif
ifeq ($(detectedOS), Windows_NT)
	JR_USER_DIR="$(LOCALAPPDATA)" //@TODO
endif
else
	JR_USER_DIR=$(XDG_DATA_HOME)
endif

hello:
	@echo "JR,the JSON Random Generator"
	@echo " Version: $(VERSION)"
	@echo " Go Version: $(GOVERSION)"
	@echo " Build User: $(USER)"
	@echo " Build Time: $(TIME)"
	@echo " Detected OS: $(detectedOS)"
	@echo " JR System Dir: $(JR_SYSTEM_DIR)"
	@echo " JR User Dir: $(JR_USER_DIR)"

install-gogen:
	go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
	#go install github.com/hamba/avro/v2/cmd/avrogen@latest

generate:
	go generate pkg/generator/generate.go

compile:
	@echo "Compiling"
	go build -v -ldflags="-s -w -X 'github.com/jrnd-io/jr/pkg/cmd.Version=$(VERSION)' \
	-X 'github.com/jrnd-io/jr/pkg/cmd.GoVersion=$(GOVERSION)' \
	-X 'github.com/jrnd-io/jr/pkg/cmd.BuildUser=$(USER)' \
	-X 'github.com/jrnd-io/jr/pkg/cmd.BuildTime=$(TIME)'" \
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
	golangci-lint run --config .localci/lint/golangci.yml --out-format tab

help: hello
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}all${RESET}'
	@echo ''

copy_templates:
	mkdir -p $(JR_SYSTEM_DIR)/$(JR_HOME)/kafka && \
	cp -r templates $(JR_SYSTEM_DIR)/$(JR_HOME) && \
	cp -r pkg/producers/kafka/*.properties.example $(JR_SYSTEM_DIR)/$(JR_HOME)/kafka/

copy_config:
	mkdir -p $(JR_SYSTEM_DIR)/$(JR_HOME) && \
	cp config/* $(JR_SYSTEM_DIR)/$(JR_HOME)/

install:
	install build/jr /usr/local/bin

all: hello install-gogen generate compile
all_offline: hello generate compile