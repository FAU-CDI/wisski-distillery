#!/bin/bash
set -e

echo "=> Building executable"
CGO_ENABLED=0 go build -o ./wdcli ./cmd/wdcli
echo "=> Running executable"
sudo "./wdcli" "$@"


# TODO: len(args)
# - remove description functions