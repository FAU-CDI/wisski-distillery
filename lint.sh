#!/bin/bash
set -e

echo "=> go vet"
go vet ./...

echo "=> staticcheck"
go tool honnef.co/go/tools/cmd/staticcheck ./...

echo "=> golangci-lint"
go tool golangci-lint run ./...

echo "=> govulncheck"
go tool govulncheck

#echo "=> gosec"
#go tool gosec ./...