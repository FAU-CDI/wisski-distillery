#!/bin/bash
set -e

# for vscode, use a debug config like the following
<<'END_DEBUG'
{
    "name": "debug.sh",
    "type": "go",
    "request": "attach",
    "mode": "remote",
    "port": 2345,
    "host": "127.0.0.1"
}
END_DEBUG

DLV=`which dlv`

echo "=> Building executable"
CGO_ENABLED=0 go build -o wdcli ./cmd/wdcli
echo "=> Running executable, attch with dlv to start"
sudo $DLV exec --only-same-user=false --headless --listen=127.0.0.1:2345 -- ./wdcli "$@"