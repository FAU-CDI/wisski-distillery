.PHONY: clean all deps

all: wdcli frontend

wdcli: internal/component/static/dist
	go build -o ./wdcli ./cmd/wdcli

internal/component/static/dist: internal/component/static/src
	rm -rf internal/component/static/dist
	cd internal/component/static/ && yarn dist

deps:
	cd internal/component/static/ && yarn install

clean:
	rm wdcli