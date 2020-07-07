# WissKI-Distillery

WissKI-Distillery is a Docker-based server provisioning and managing for multiple [WissKI](https://wiss-ki.eu/) instances. 

**This project is still a work in progress and nothing in this repository is ready for production use** 

## Overview

This project consists of the following:

- this README
- bash scripts for setting up and managing the distillery server
- bash scripts for backing up the server
- a `Vagrantfile` for local testing

The bash scripts are dependency-free and only assume that a basic debian system is available. 
The scripts have been tested only under Debian 10, but may also work under older or newer versions. 
All scripts expect to be run as root, and will fail when this is not the case. 
Each script is well-commented and all commands are explained. 

Configuration of the bash scripts can be done in the file `distillery/.env`. 
A sample configuration file (with documented defaults) is available in `distillery/.env.sample`. 
To get started, it is sufficient to run:

```bash
cd distillery/
cp .env.sample .env
your-favorite-editor .env # open and customize, usually only the domain needs adjusting
```

## Vagrantfile

For local testing, it is recommended to use [Vagrant](https://www.vagrantup.com/) and the provided `Vagrantfile`. 
After installing vagrant, run:

```bash
# start the vargant box
vagrant up

# open a shell inside the vm
# for debugging purposes forward port 7200 (GraphDB) and 8080 (phpmyadmin)
vagrant ssh -- -L 7200:127.0.0.1:7200 -L 8080:127.0.0.1:8080
```


## Preparing the Server -- 'system_install.sh'

*TLDR: `sudo bash /distillery/system_install.sh /path/to/graphdb.zip`*

To prepare the server for becoming a WissKI factory, several core Docker Instances must be installed. 
These are:

- [nginx-proxy](https://github.com/nginx-proxy/nginx-proxy) - an automated nginx reverse proxy
    - This will delegate individual hostnames to appropriate docker containers, see [this blog post](http://jasonwilder.com/blog/2014/03/25/automated-nginx-reverse-proxy-for-docker/) for an overview. 
    - Optionally makes use of [docker-letsencrypt-nginx-proxy-companion](https://github.com/nginx-proxy/docker-letsencrypt-nginx-proxy-companion) to automatically provision and renew HTTPS certificates. 
    - See [distillery/resources/compose/web](distillery/resources/compose/web) for implementation details. 

- [MariaDB](https://mariadb.org/) - an SQL server
    - It is configured to run inside a docker container
    - A passwordless `root` account is created, which can only be used from inside the container. 
    - A `bookkeeping` database and table is created by default, to store known WissKI instance metadata in. 
    - A database shell can be opened using `sudo /distillery/mysql.sh`. 
    - A [phpmyadmin](https://www.phpmyadmin.net/) is started on `127.0.0.1:8080`. 
    - See [distillery/resources/compose/sql](distillery/resources/compose/sql) for implementation details. 

- [GraphDB](http://graphdb.ontotext.com/) - a SPARQL backend for WissKI
    - It is configured to run inside a docker container. 
    - The Workbench API is started on `127.0.0.1:7200`. 
    - Security is not enabled at the moment. 
    - See [distillery/resources/compose/triplestore](distillery/resources/compose/triplestore) for implementation details. 

To manage multiple docker containers, this script makes heavy use of [docker-compose](https://docs.docker.com/compose/). 

Setting up these steps is fully automatic.
In particular, after obtaining a license and the installation zip file for 'GraphDB', one can run the 'distillery/system_install.sh' script as follows to setup all components:

```bash
sudo bash /distillery/system_install.sh /path/to/graphdb.zip
```

In principle this script is idempotent, meaning it can be run multiple times achieving the same effect. 

## Updating the Docker Containers -- 'system_update.sh'

For security purposes, the core containers should be regularly updated. 
To achieve this, the docker container images should be rebuilt and restarted. 

This can be done using:

```bash
sudo bash /distillery/system_update.sh
```

## Provisioning a new WissKI instance  -- 'provision.sh'

*TLDR: `sudo /distillery/provision.sh slug-of-new-website`*

A new WissKI instance consists of several components:

- A Drupal instance inside a lightweight php runtime container
- An entry in the SQL bookkeeping table that stores instance meta-data
- An SQL database and user for Drupal
- A GraphDB repository and user as SPARQL endpoint

Each WissKI instance is identified by a ``slug''. 
This is a preferably short name that is used to form a domain name for the WissKI instance. 
The WissKI distillery assumes that each instance is a subdomain of a given domain. 
For example, if the given domain is 'wisskis.example.com' and the slug of a particular instance is 'blue', the subdomain used by this instance would be 'blue.wisskis.example.com'. 
The given domain can be configured within the '.env' file. 

We use the following process to provision a new instance:

__1. Create a new docker-compose.yml file__

In this step we first create a directory on the real system to hold all files relating to this instance. 
By default, this takes place inside `/var/www/deploy/instances/$DOMAIN`, but this can be configured. 
We then create a docker-compose file in this directory that is ready for running the runtime container. 

__2. Create an appropriate SQL database and user__

We create a new SQL database to eventually store Drupal-related data in. 
The user and database names are generated from the slug. 
The database password is randomly generated and only made available directly to the Drupal instance later. 

__3. Create a GraphDB repository and user__

Next, we create a dedidcated GraphDB repository for the WissKI instance. 
We also create a new GraphDB user with access to this repository. 

__4. Provision the instance inside the container__

We start the container in provisioning mode. 

This does the following:

- Creates a new composer project that requires [drupal/recommended-project](https://github.com/drupal/recommended-project)`. 
- Installs `drush` into this project. 
- Runs the `drush site-install` command to configure the Drupal instance. Generates a random password to use. 
- Adds and enables WissKI-specific modules for this instance. 

Currently the WissKI Salz instance is not enabled programatically. 
Instead all credentials (along with instructions on how to configure it) are printed to the command line. 


__6. Start the Docker Container__

Finally, we can start the docker container. 

These steps can be performed automatically. 
To do so, use:

```bash
sudo bash /distillery/provision.sh SLUG
```

## Purge an existing WissKI instance -- 'purge.sh'


Sometimes it is required to remove a given WissKI instance. 
In particular all parts belonging to it should be removed. 

To use it, run:

```bash
sudo bash /distillery/purge.sh SLUG
```

## Open a shell -- 'shell.sh'

Sometimes manual changes to a given WissKI instance are required. 
For this purpose, you can use:

```bash
sudo bash /distillery/shell.sh SLUG
```

This will open a shell in the provided WissKI instance. 

## License

This project and associated files in this repository are licensed as follows:

    WissKI-Distillery - A docker-based WissKI instance server
    Copyright (C) 2020 AGFD <https://www.agfd.fau.de/>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

Please see `LICENSE` for a legally binding license text. 
The short summary of the license is:

- You may use this software for any purpose, including commerical. 
- You may create derivative works, and use those for any purpose, including commerical. 

if you follow the following conditions:

- You provide the end-user with a copy of this license. 
- You make the source code of any derivative works available. 
- Any derivative works clearly list changes made. 
- You license any derivative works under the same license. 

This also applies if you only run a backend service based on this software. 


## TODO

- User-level documentation
    - What is a factory?
    - Why a factory?
    - First steps after provisioning
- Automatically setup SALZ adapter (if this is possible)
- Enable authentication for GraphDB
- Investigate support for GraphDB Auth in WissKI Salz
    - Eventually enable security if needed
    - Switch to a different TripleStore alltogether?
- Investigate managing phpmyadmin
- Investigate managing graphdb
- Investigate delegating shell access
- Investigate delegating ftp access
- document CNAME structure