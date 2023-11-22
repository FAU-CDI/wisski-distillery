.PHONY: clean all deps live tslint tsfix

live:
	sudo CGO_ENABLED=0 go run -trimpath ./cmd/wdcli $(ARGS)

all: wdcli

wdcli:
	go generate ./internal/dis/component/control/static/
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