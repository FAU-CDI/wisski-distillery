# Automatic Drupal and WissKi factory scripts

This repository contains a factory server implementation that creates and maintains a list of Drupal Instances. 

** This is a work in progress and nothing in this repository is ready for production use ** 

## Overview

This project consists of the following:

- this README
- several bash scripts in the 'factory' folder that are described below
- a `Vagrantfile` for local testing

The bash scripts are dependency-free and only assume that a basic debian system is available. 
The scripts have been tested only under Debian 10, but may also work under older or newer versions. 
All scripts expect to be run as root, and will fail when this is not the case. 
Each script is well-commented and all commands are explained. 

Configuration of the bash scripts can be done in the file 'factory/.env'. 
A sample configuration file (with documented defaults) is available in 'factory/.env.sample'. 
To get started, it is sufficient to run:

```bash
cd factory/
cp .env.sample .env
your-favorite-editor .env # open and customize, usually only the domain needs adjusting
```

For local testing, it is recommended to use [Vagrant](https://www.vagrantup.com/) and the provided `Vagrantfile`. 

## Preparing the Server -- 'system_install.sh'

*TLDR: `sudo bash /factory/system_install.sh /path/to/graphdb.zip`*

To prepare the server for becoming a WissKI factory, a few components need to be installed. 
In particular, these are:
- [PHP](https://www.php.net/) and [Composer](https://getcomposer.org/) -- for getting and running the Drupal Code
- serveral PHP modules that are dependencies of Drupal
- [MariaDB](https://mariadb.org/) -- an SQL database
- [Apache2](https://httpd.apache.org/), the corresponding php and [mpm-itk](http://mpm-itk.sesse.net/) modules -- a webserver
- [GraphDB](http://graphdb.ontotext.com/) - an SPARQL backend for WissKi

With the exception of GraphDB all these components can be installed using Debian's package manager 'apt'. 
To install GraphDB, a zip with the binaries needs to be unpacked, and then a systemd service for it needs to be created. 

These steps can be performed automatically. 
In particular, after obtaining a license and the installation zip file for 'GraphDB', one can run the 'factory/system_install.sh' script as follows to setup all components:

```bash
sudo bash /factory/system_install.sh /path/to/graphdb.zip
```

In principle this script is idempotent, meaning it can be run multiple times achieving the same effect. 

## Provisioning a new WissKi instance  -- 'provision.sh'

*TLDR: `sudo bash provision.sh slug-of-new-website`*

A new WissKi instance consists of several components:

- A [Drupal](https://www.drupal.org/) instance, managed as a [Composer](https://getcomposer.org/) project
- An [Apache](https://httpd.apache.org/) the makes the above available externally
- An [SQL](https://mariadb.org/) database, to store Drupal Nodes in
- A [GraphDB](https://graphdb.ontotext.com/) repository to store RDF triples in

Each WissKi instance is identified by a ``slug''. 
This is a preferably short name that is used to form a domain name for the WissKi instance. 
This factory assumes that each instance is a subdomain of a given domain. 
For example, if the given domain is 'wisskis.example.com' and the slug of a particular instance is 'blue', the subdomain used by this instance would be 'blue.wisskis.example.com'. 
The given domain can be configured within the '.env' file. 

In this implementation we furthermore isolate each WissKi instance from the rest of the system.
For this purpose, we make use of an appropriate system user, an appropriate SQL user and a GraphDB user. 
**Note: GraphDB users are not yet implemented **

We thus use the following process to provision a new instance:

__1. We create a new system user and hoem directory__

The username is derived from the slug, with a configurable prefix. 
The home directory for this user will contain the Drupal PHP files needed to run a WissKi. 
For this reason, the home directory for each user is a subdirectory at a standardized location. 
By default this is `/var/www/factory/$USER', but this can be customized. 

__2. Create an appropriate SQL database and user__

We create a new SQL database to eventually store Drupal-related data in. 
The user and database names are again generated from the slug. 
The database password is randomly generated and only made available directly to the Drupal instance later. 

__3. Initialize a new composer project__

Within the home directory of the dedicated user, we create a new composer project that requires [drupal/recommended-project](https://github.com/drupal/recommended-project)` as well as drush. 

__4. Run the Drupal Installation scripts__

We run the Drupal installation scripts. 
Here we tell Drupal about the database credentials, and initialize an initial 'admin' user for the drupal instance. 
The password for the 'admin' user is randomly generated in this process. 

__5. Create a GraphDB repository__

Next, we create a dedidcated GraphDB repository for the WissKi instance. 
*TODO*: Create a GraphDB user. 

__6. Add WissKi modules to Drupal__

Next, we add the required WissKi modules to Drupal. 
*TODO*: Configure the WissKi modules automatically. 

__7. Create a Apache VHost configuration__

Finally, we create an apache vhost configuration that makes the drupal website available. 
*TODO*: SSL


These steps can be performed automatically. 
To do so, use:

```bash
sudo bash /factory/provision.sh SLUG
```

## Manually editing WissKi instances -- 'shell.sh'

Sometimes it is needed to make manual adjustments to an individual instance. 
For this purpose, the `shell.sh` script exists. 
It opens an interactive shell in the context of a given WissKi instance. 
In particular it:
- switches to the appropriate system user
- sets up the '$PATH' environment variable to allow using 'drush' and 'composer'

To use it, run:

```bash
sudo bash /factory/shell.sh SLUG
```

## Removing an existing WissKi instance -- 'remove.sh'

* TODO: Document this more *


Sometimes it is required to remove a given WissKi instance. 
In particular all parts belonging to it should be removed. 


To use it, run:

```bash
sudo bash /factory/remove.sh SLUG
```


## TODO

- More documentation
    - Document and improve`update.sh`
    - User-level documentation
        - What is a factory?
        - Why a factory?
        - First steps after provisioning
- Writeup approach to SSL (Wildcard cert with proxy that downgrades connections to plain http, or mod_md)
- Automatically setup SALZ adapter (if this is possible)
- Setup users for GraphDB and enable security, is this supported by WissKi SALZ?
- Allow customization of GraphDB paths


## License

Licensed under GPL 3. 