#!/bin/bash
set -e

echo "=> Creating 'sites/default/libraries/dompurity/dist/' directory"
mkdir -p /var/www/data/project/web/libraries/dompurity/dist/

echo "=> Downloading 'purify.min.js' and 'LICENSE'"
curl -L https://raw.githubusercontent.com/cure53/DOMPurify/main/dist/purify.min.js -o /var/www/data/project/web/libraries/DOMPurify/dist/purify.min.js
curl -L https://raw.githubusercontent.com/cure53/DOMPurify/main/LICENSE -o /var/www/data/project/web/libraries/DOMPurify/LICENSE

echo "=> Done"