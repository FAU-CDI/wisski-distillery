#!/bin/bash

# This script will provision a new Drupal instance and make it available to apache. 
# Usage: sudo ./provision.sh $SLUG
# In case the installation fails, it will bail out and leave you with an incomplete installation. 
# To delete an incomplete installation, use the ./remove.sh script, or try fixing the error manually. 
set -e

# read the lib/shared.sh and read the slug argument. 
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
require_slug_argument


# Check that the apache2 config is correct. 
# This is a sanity test so that we don't randomly fail later because of bad config. 
log_info " => Checking apache configuration"
apache2ctl configtest > /dev/null

# Create a system user and group
log_info " => Creating system user and group '$SYSTEM_USER'"
addgroup --system "$SYSTEM_USER"
adduser --home "$BASE_DIR" --system --disabled-password --disabled-login --ingroup "$SYSTEM_USER" "$SYSTEM_USER"

# Make directory for the composer project to live in
log_info " => Making composer directory '$COMPOSER_DIR'"
sudo -u "$SYSTEM_USER" mkdir -p "$COMPOSER_DIR"
cd "$COMPOSER_DIR"

# Write out a new apache configuration file into /etc/apache2/sites-available. 
# We will need to substiute in some configuration directories. 
log_info " => Writing new apache configuration file"
cat << EOF >> "$APACHE_CONFIG_SITE_AVAILABLE"
<VirtualHost *:80>
    DocumentRoot $WEB_DIR
    ServerName $INSTANCE_DOMAIN
    AssignUserId $SYSTEM_USER $SYSTEM_USER

    <Directory $WEB_DIR>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>
    ErrorLog \${APACHE_LOG_DIR}/error.log
    CustomLog \${APACHE_LOG_DIR}/access.log combined
</VirtualHost>
EOF

# Create a new composer project. 
log_info " => Creating composer project"
composer create-project drupal/recommended-project .
composer require drush/drush

# Randomly generate the database name and user we will configure. 
# Use the 'randompw' alias for this. 
log_info " => Generating new MySQL password"
MYSQL_PASSWORD="$(randompw)"

# Initialize the SQL database with those credentials. 
log_info " => Intializing new SQL database '${MYSQL_DATABASE}' and user '$MYSQL_USER'. "
mysql -e "CREATE DATABASE \`${MYSQL_DATABASE}\`;"
mysql -e "CREATE USER \`${MYSQL_USER}\`@localhost IDENTIFIED BY '${MYSQL_PASSWORD}';"
mysql -e "GRANT ALL PRIVILEGES ON \`${MYSQL_DATABASE}\`.* TO \`${MYSQL_USER}\`@localhost;"
mysql -e "FLUSH PRIVILEGES;"

# Generate some more random credentials, this time for drupal. 
# We again make use of the randompw alias. 
log_info " => Generating new drupal credentials"
DRUPAL_USER="admin"
DRUPAL_PASS="$(randompw)"

# Use 'drush' to run the site-installation. 
# Here we need to use the username, password and database creds we made above. 
log_info " => Running drupal installation scripts"
drush site-install standard --yes --site-name=${INSTANCE_DOMAIN} --account-name=$DRUPAL_USER --account-pass=$DRUPAL_PASS --db-url=mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@localhost/${MYSQL_DATABASE}
drupal_sites_permission_workaround

# Create a new repository for GraphDB. 
# First write out the configuration into a new directory. 
log_info " => Writing GraphDB configuration in temporary directory"
tmpdir="$(mktemp -d)"
cd "$tmpdir"
cat << EOF > repo-config.ttl
# Creates a new GraphDB repository with Wisski

@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#>.
@prefix rep: <http://www.openrdf.org/config/repository#>.
@prefix sr: <http://www.openrdf.org/config/repository/sail#>.
@prefix sail: <http://www.openrdf.org/config/sail#>.
@prefix owlim: <http://www.ontotext.com/trree/owlim#>.

[] a rep:Repository ;
    rep:repositoryID "$GRAPHDB_REPO" ;
    rdfs:label "$INSTANCE_DOMAIN" ;
    rep:repositoryImpl [
        rep:repositoryType "graphdb:FreeSailRepository" ;
        sr:sailImpl [
            sail:sailType "graphdb:FreeSail" ;

            owlim:owlim-license "" ;

            owlim:base-URL "http://$INSTANCE_DOMAIN/owlim#" ;
            owlim:defaultNS "" ;
            owlim:entity-index-size "10000000" ;
            owlim:entity-id-size  "32" ;
            owlim:imports "" ;
            owlim:repository-type "file-repository" ;
            owlim:ruleset "empty" ;
            owlim:storage-folder "storage" ;

            owlim:enable-context-index "false" ;
            owlim:cache-memory "80m" ;
            owlim:tuple-index-memory "80m" ;

            owlim:enablePredicateList "false" ;
            owlim:predicate-memory "0%" ;

            owlim:fts-memory "0%" ;
            owlim:ftsIndexPolicy "never" ;
            owlim:ftsLiteralsOnly "true" ;

            owlim:in-memory-literal-properties "false" ;
            owlim:enable-literal-index "true" ;
            owlim:index-compression-ratio "-1" ;

            owlim:check-for-inconsistencies "false" ;
            owlim:disable-sameAs  "false" ;
            owlim:enable-optimization  "true" ;
            owlim:transaction-mode "safe" ;
            owlim:transaction-isolation "true" ;
            owlim:query-timeout  "0" ;
            owlim:query-limit-results  "0" ;
            owlim:throw-QueryEvaluationException-on-timeout "false" ;
            owlim:useShutdownHooks  "true" ;
            owlim:read-only "false" ;
            owlim:nonInterpretablePredicates "http://www.w3.org/2000/01/rdf-schema#label;http://www.w3.org/1999/02/22-rdf-syntax-ns#type;http://www.ontotext.com/owlim/ces#gazetteerConfig;http://www.ontotext.com/owlim/ces#metadataConfig" ;
        ]
    ].
EOF

# Create the configuration and use the configuration generated above. 
# TODO: Permissions for GraphdDB
log_info "Generating new GraphDB repository '$GRAPHDB_REPO'"
curl -X POST\
    http://127.0.0.1:7200/rest/repositories\
    -H 'Content-Type: multipart/form-data'\
    -F "config=@repo-config.ttl"

# Remove the temporary directory. 
cd ..
rm -rf "$tmpdir"

# Install the Wisski packages. 
log_info " => Installing Wisski packages"
cd "$COMPOSER_DIR"

drupal_sites_permission_workaround
composer require drupal/wisski

drupal_sites_permission_workaround
composer require drupal/inline_entity_form

drupal_sites_permission_workaround
composer require drupal/imagemagick

drupal_sites_permission_workaround
composer require drupal/image_effects

drupal_sites_permission_workaround
composer require drupal/colorbox


# Enable the WissKi modules. 
log_info " => Enable Wisski modules"
drush pm-enable --yes wisski_core wisski_linkblock wisski_pathbuilder wisski_adapter_sparql11_pb wisski_salz
drupal_sites_permission_workaround

# TODO: Setup WissKi-Salz. 

# Finally enable the apache2 config. 
# And then reload to start serving it. 
log_info " => Enabling and reloading apache configuration"
a2ensite "${INSTANCE_DOMAIN}"
systemctl reload apache2

# TODO: Certbot support

# and done!
log_info " => Finished, your Drupal details are: "
echo "URL:                  http://$INSTANCE_DOMAIN"
echo "Username:             $DRUPAL_USER"
echo "Password:             $DRUPAL_PASS"
log_info " => Your GraphDB details (for WissKi Salz) are: "
echo "Read URL:             http://127.0.0.1:7200/repositories/$GRAPHDB_REPO"
echo "Write URL:            http://127.0.0.1:7200/repositories/$GRAPHDB_REPO/statements"
echo "Writable:             yes"
echo "Default Graph URI:    http://$INSTANCE_DOMAIN/owlim#"
echo "Ontology Paths:       (empty)"
echo "SameAs property:      http://www.w3.org/2002/07/owl#sameAs"