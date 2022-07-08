# Init variables
VERSION ?= $(shell git describe --tags --always)
COMMIT ?= $(shell git rev-parse HEAD)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = "-w -X main.Version=$(VERSION) -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

GITHUB_ORG_NAME = bjornm82
GITHUB_PROJECT = trading-obs
PROJECT_DIR ?= ${GOPATH}/src/github.com/${GITHUB_ORG_NAME}/${GITHUB_PROJECT}

APP ?= go-trading-obs
OS ?= linux
ARCH ?= amd64

.PHONY: all
all:
	$(MAKE) deps
	$(MAKE) build
	$(MAKE) run-orderer
	$(MAKE) run-positioner

.PHONY: deps
deps:
	go mod tidy

.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(APP) -a -tags netgo -ldflags $(LDFLAGS) .

.PHONY: run-orderer
run-orderer:
	docker-compose up --build orderer

.PHONY: run-positioner
run-positioner:
	docker-compose up --build trading_obs

.PHONY: test
test:
	go test -race -v ./... | grep -v vendor