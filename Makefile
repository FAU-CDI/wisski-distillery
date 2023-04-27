.PHONY: clean all deps live

live:
	sudo CGO_ENABLED=0 go run ./cmd/wdcli $(ARGS)

all: wdcli

wdcli:
	go generate ./internal/dis/component/control/static/
	CGO_ENABLED=0 go build -o ./wdcli ./cmd/wdcli

deps: internal/dis/component/server/assets/node_modules

internal/dis/component/server/assets/node_modules:
	cd internal/dis/component/server/assets/ && yarn install

clean:
	rm wdcli