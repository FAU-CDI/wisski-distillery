.PHONY: clean all deps

all: wdcli

wdcli:
	go generate ./internal/dis/component/control/static/
	go build -o ./wdcli ./cmd/wdcli

deps:
	cd internal/dis/component/control/static/ && yarn install

clean:
	rm wdcli