.PHONY: clean all

all: wdcli

wdcli:
	go build -o ./wdcli ./cmd/wdcli

clean:
	rm wdcli