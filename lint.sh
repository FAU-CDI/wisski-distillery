#!/bin/bash
set -e

echo "=> go vet"
go vet ./...

echo "=> golangci-lint"
go tool golangci-lint run ./...

echo "=> govulncheck"

#echo "=> gosec"
go tool gosec ./...