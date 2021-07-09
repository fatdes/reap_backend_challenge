# Simple website

## Setup

- golang 1.16+
- make
- docker
- [gauge](https://docs.gauge.org/getting_started/installing-gauge.html?os=macos&language=javascript&ide=vscode)

```bash
# for generating unit tests mocks
## install mockgen
go install github.com/golang/mock/mockgen@v1.6.0
## setup path to the installed mockgen if needed
export PATH=$PATH:$(go env GOPATH)/bin
```

```bash
# for automated tests
## install gauge
https://docs.gauge.org/getting_started/installing-gauge.html
```

## Documentation

1. Make command to test
```bash
# generate mocks
make generate

# test
make test

# run automated test against local web server
make local-automated-test

# clean up
make local-clean
```

2. API documentation [swagger doc format](api.yaml)