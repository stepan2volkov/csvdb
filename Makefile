check:
	golangci-lint run -c golangci-lint.yaml

test:
	go test ./...