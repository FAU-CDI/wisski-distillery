.PHONY: clean all deps

all: wdcli

wdcli:
	go generate ./internal/dis/component/static/
	go build -o ./wdcli ./cmd/wdcli

deps:
	cd internal/component/static/ && yarn install

clean:
	rm wdcli