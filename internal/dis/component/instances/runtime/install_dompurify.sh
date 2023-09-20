#!/bin/bash
set -e

mkdir -p /var/www/data/project/web/libraries/DOMPurify/dist/
curl -L https://raw.githubusercontent.com/cure53/DOMPurify/main/dist/purify.min.js -o /var/www/data/project/web/libraries/DOMPurify/dist/purify.min.js