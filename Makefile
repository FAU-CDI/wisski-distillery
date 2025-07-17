.PHONY: clean all deps live tslint tsfix lint

all: wdcli

lint:
	go vet ./...
	go tool govulncheck ./...
	go tool golangci-lint run ./...

wdcli:
	go generate ./internal/dis/component/server/assets/
	CGO_ENABLED=0 go build -trimpath -o ./wdcli ./cmd/wdcli

tslint:
	cd internal/dis/component/server/assets/ && yarn ts-standard

tsfix:
	cd internal/dis/component/server/assets/ && yarn ts-standard --fix

deps: internal/dis/component/server/assets/node_modules

internal/dis/component/server/assets/node_modules:
	cd internal/dis/component/server/assets/ && yarn install

clean:
	rm wdcli