#!/bin/bash
set -e

# read the lib/shared.sh
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"

# This script will prepare a server to become a factory for Drupal Instances. 
# Even though it assumes a clean server, it *should* be idempotent. 
log_info "=> Preparing system to serve as a factory. "

# Read the 'GRAPHDB_ZIP' argument from the command line. 
# If it's not set, throw an error. 
GRAPHDB_ZIP=$1
if [ -z "$GRAPHDB_ZIP" ]; then
    log_error "Usage: system_install.sh GRAPHDB_ZIP"
    exit 1;
fi;

# Make a temporary directory to use for various tasks during this script. 
log_info " => Making temporary directory"
tmpdir="$(mktemp -d)"
log_ok "Made $tmpdir"

# fetch new package versions, then upgrade everything we already have. 
# This isn't technically neccessary, but it means it'll work on an otherwise untouched system. 
log_info " => Installing package updates ..."
apt-get update
apt-get dist-upgrade -y

# Install composer, by downloading it using curl and then run it with php to install it in /usr/local/bin. 
log_info " => Installing composer"
apt-get install -y curl php-cli php-mbstring git unzip
curl -sS https://getcomposer.org/installer -o "$tmpdir/composer-setup.php"
php $tmpdir/composer-setup.php --install-dir=/usr/local/bin --filename=composer

# Install required php extensions for Drupal and WissKi. 
log_info " => Installing required php extensions"
apt-get install -y php-xml php-gd php-mysql php-common php-xmlrpc php-soap php-gd php-intl php-mysql php-zip php-curl php-ssh2

# Install the mariadb kernel. 
log_info " => Installing mariadb"
apt-get -y install mariadb-server

# Install apache and required php extensions. 
log_info " => Installing apache2, php and auth modules"
apt-get install -y apache2 libapache2-mod-php libapache2-mpm-itk

# Make the directory for all drupal instances to live in. 
log_info " => Making root directory for Drupal Installations"
mkdir -p "$DRUPAL_ROOT"

# Install java for GraphDB. 
# We use the 'headless' package to prevent installing anything graphical on a headless server. 
log_info " => Installing java"
apt-get install -y default-jre-headless


# Next we have to check if we need to install graphdb. 
# If '/opt/graphdb' exists, assume that the installation has already been performed. 
if [ -d "/opt/graphdb" ]; then
    log_info " => 'opt/graphdb' exists, skipping setup step. ";
else

    # Unzip the GraphDB sources into a temporary directory. 
    echo " => Unzipping GraphDB into temporary directory"
    unzip "$GRAPHDB_ZIP" -d "$tmpdir/graphdb"

    # Then move them into /opt/graphdb.
    # Here we need to make sure that the first subdirectoy is renamed appropriately during the move. 
    echo " => Moving GraphDB into /opt/graphdb"
    mv "$tmpdir/graphdb"/* /opt/graphdb
fi

# Next make a system group 'graphdb' and system user 'graphdb'. 
# And also chown the /opt/graphdb directory to that user. 
# As the user might already exist, we surpress errors of the commands. 
log_info " => Making GraphDB group and user"
addgroup --system graphdb || true
adduser --home "/opt/graphdb" --system --no-create-home --disabled-password --disabled-login --ingroup graphdb graphdb || true
chown -R graphdb:graphdb /opt/graphdb

# Create a service file to use graphdb with systemd. 
# This file uses the users created above, and also hard-codes listening address and maximum memory. 
# This avoids having to write the config file using bash hacks. 
log_info " => Making 'graphdb.service'"
cat << "EOF" > /etc/systemd/system/graphdb.service
[Unit]
Description=GraphDB

[Service]
Type=simple
User=graphdb
Group=graphdb
ExecStart=/opt/graphdb/bin/graphdb –Xmx6g -Dgraphdb.connector.address=127.0.0.1

[Install]
WantedBy=multi-user.target
EOF

# We just created a service, so now start it and put it into autostart mode. 
log_info " => Starting and enabling graphdb.service"
systemctl enable graphdb
systemctl start graphdb

# Finally remove the temporary directory we created above. 
log_info " => Removing temporary directory"
rm -rf "$tmpdir"

log_info " => Server is now ready to become a factory. "