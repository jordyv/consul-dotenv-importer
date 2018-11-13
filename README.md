# Consul dotenv file importer #

Simple CLI tool written in Go to import a .env file in Consul.

## Installation ##

```
go get -u github.com/jordyv/consul-dotenv-import
```

## Usage ##

```
$ consul-dotenv-import -h
Usage: consul-dotenv-import <options> <env file>
  -dry
        Dry run (don't actual put values in Consul)
  -host string
        Consul host (default "http://127.0.0.1:8500")
  -prefix string
        Consul KV prefix
  -token string
        Consul token
```

### Example ###

Env file:
```dotenv
TEST=foo
BAR=test
```

Run:
```
$ consul-dotenv-import -prefix testenv .env
Putting testenv/TEST = foo
Putting testenv/BAR = test
```

## Development ##

```
$ git clone https://github.com/jordyv/consul-dotenv-import
$ cd consul-dotenv-import
$ dep ensure
$ go run main.go
```
