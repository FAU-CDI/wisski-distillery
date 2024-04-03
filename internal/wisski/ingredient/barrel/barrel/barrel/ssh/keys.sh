#!/bin/bash

curl -H "Host:$(hostname -f)" http://dis:9999/authorized_keys
