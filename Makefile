BUILD_COMMIT := $(shell git log --format="%H" -n 1)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%S')

PROJECT = github.com/stepan2volkov/csvdb
ENTRYPOINT = cmd/csvdb
CMD:= $(PROJECT)/$(ENTRYPOINT)

check:
	./bin/golangci-lint run -c golangci-lint.yaml

test:
	go test -cover ./...

run:
	go run $(ENTRYPOINT)/csvdb.go

.PHONY: build
build:
	mkdir -p build
	go build -ldflags="\
		-X '$(PROJECT)/internal/app.BuildCommit=$(BUILD_COMMIT)'\
		-X '${PROJECT}/internal/app.BuildTime=${BUILD_TIME}'"\
		-o build $(CMD)

clean:
	rm -rf build

install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.45.2
