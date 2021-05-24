default: all

SRC = getopt.go getopt_test.go

all: build test lint vet

build: $(SRC)
	go build

test: $(SRC)
	go test

lint: $(SRC)
	golint

vet: $(SRC)
	go vet
