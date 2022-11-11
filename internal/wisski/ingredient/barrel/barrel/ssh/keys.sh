#!/bin/bash

cat /var/www/.ssh/authorized_keys /var/www/.ssh/global_authorized_keys 2> /dev/null || exit 0
