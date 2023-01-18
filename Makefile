.PHONY: clean all deps

all: wdcli

wdcli:
	go generate ./internal/dis/component/control/static/
	go build -o ./wdcli ./cmd/wdcli

deps: internal/dis/component/control/static/node_modules

internal/dis/component/control/static/node_modules:
	cd internal/dis/component/control/static/ && yarn install

clean:
	rm wdcli