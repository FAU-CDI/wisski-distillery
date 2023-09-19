#!/bin/bash
set -e

# make the target directory and install
mkdir -p /var/www/data/project/web/libraries/DOMPurify/dist/
curl -L https://raw.githubusercontent.com/cure53/DOMPurify/main/dist/purify.min.js -o /var/www/data/project/web/libraries/DOMPurify/dist/purify.min.js