.PHONY: build install

build:
	go build -o dist/consul-dotenv-import main.go

install:
	go install
