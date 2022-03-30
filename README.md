# WissKI-Distillery

WissKI-Distillery is a Docker-based server provisioning and managing for multiple
[WissKI](https://wiss-ki.eu/) instances.

The WissKI Distillery is a set of scripts, tools, and applications that allows to operate
a WissKI cloud of distinct but jointly managed WissKI instances, hosted on a dedicated
hardware pool. Like the WissKI system, the WissKI Distillery is open source and free to
use.

This README contains only technical documentation.
For members of [FAU Erlangen-Nürnberg](https://www.fau.de/) a cloud offering based on this service known as FAUWissKICloud.
Please see https://wisski.agfd.fau.de/ for related documentation. 

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
# once, to install the plugin to automatically build the guest iso
# at the time of writing version 0.25.0 is broken.
vagrant plugin install --plugin-version 0.24.0 vagrant-vbguest

# start the vargant box
vagrant up

# open a shell inside the vm
# for debugging purposes forward port 7200 (GraphDB) and 8080 (phpmyadmin)
vagrant ssh -- -L 7200:127.0.0.1:7200 -L 8080:127.0.0.1:8080
```

## Preparing the Server -- 'system_install.sh'

_TLDR: `sudo bash /distillery/system_install.sh /path/to/graphdb.zip`_

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

- [proxyssh](https://github.com/tkw1536/proxyssh) - an ssh server that delegates client connections to different WissKIs
  - It is configured to run inside a docker container
  - Uses a global configurable authorized_keys file.
  - Also allows users to write their own authorized_keys files.

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

## Provisioning a new WissKI instance -- 'provision.sh'

_TLDR: `sudo /distillery/provision.sh slug-of-new-website`_

A new WissKI instance consists of several components:

- A Drupal instance inside a lightweight php runtime container (a `barrel` in which to store WissKI)
- An entry in the SQL bookkeeping table that stores instance meta-data
- An SQL database and user for Drupal
- A GraphDB repository and user as SPARQL endpoint

Each WissKI instance is identified by a ``slug''.
This is a preferably short name that is used to form a domain name for the WissKI instance.
The WissKI distillery assumes that each instance is a subdomain of a given domain.
For example, if the given domain is 'wisskis.example.com' and the slug of a particular instance is 'blue', the subdomain used by this instance would be 'blue.wisskis.example.com'.
The given domain can be configured within the '.env' file.

We use the following process to provision a new instance:

**1. Create a new docker-compose.yml file**

In this step we first create a directory on the real system to hold all files relating to this instance.
By default, this takes place inside `/var/www/deploy/instances/$DOMAIN`, but this can be configured.
We then create a docker-compose file in this directory that is ready for running the `barrel` container.

**2. Create an appropriate SQL database and user**

We create a new SQL database to eventually store Drupal-related data in.
The user and database names are generated from the slug.
The database password is randomly generated and only made available directly to the Drupal instance later.

**3. Create a GraphDB repository and user**

Next, we create a dedicated GraphDB repository for the WissKI instance.
We also create a new GraphDB user with access to this repository.

**4. Provision the instance inside the container**

We start the container in provisioning mode.

This does the following:

- Creates a new composer project that requires [drupal/recommended-project](https://github.com/drupal/recommended-project)`.
- Installs `drush` into this project.
- Runs the `drush site-install` command to configure the Drupal instance. Generates a random password to use.
- Adds and enables WissKI-specific modules for this instance.

Currently the WissKI Salz instance is not enabled programatically.
Instead all credentials (along with instructions on how to configure it) are printed to the command line.

**6. Start the Docker Container**

Finally, we can start the docker container.

These steps can be performed automatically.
To do so, use:

```bash
sudo bash /distillery/provision.sh SLUG
```

## Rebuild an instance -- 'rebuild.sh' and 'rebuild-all.sh'

Sometimes it becomes necessary (because of changes to this project) to rebuild the docker image running a certain docker instance.
To do so, use:

```bash
sudo bash /distillery/rebuild.sh SLUG
```

Note that rebuilding an instance does restart the docker container resulting in a small (typical < 1 second) interruption to the website in question.
Furthermore, while the container recreated, the old image stays on the host.
To delete all instances, run:

```bash
sudo docker image prune --all
```

To automatically rebuild all instances, use:

```bash
sudo bash /distillery/rebuild-all.sh
```

## Reserving an instance -- 'reserve.sh'

Sometimes it is useful to reserve a particular instance name.
This is done by hosting a placeholder website at the domain.
To do so, use:

```bash
sudo bash /distillery/reserve.sh SLUG
```

To un-reserve a website, manually stop the docker stack and remove the folder.

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

## List all instances -- 'ls.sh'

To list all instances, the following command can be used:

```bash
sudo bash /distillery/ls.sh
```

## Backups -- 'backup_all.sh' and 'backup_instance.sh'

This project comes with a backup script.
To make a backup of *all instances*, run:

```bash
sudo bash /distillery/backup_all.sh
```

Backups may temporarily shutdown individual instances to ensure data consistency.
Typical backup times are a minute or less.

Backups are stored in the `backups/final` directory.
They contain:

- a backup of every single instance (see below)
- a complete backup of the SQL database
- nquads of all the GraphDB repositories
- a backup of the config file

Files are `.tar.gz`ipped.
By default, backups are kept for up to thirty days, after which they are removed.

This script does not automatically provision a cronjob.
An example job to e.g. run a backup every saturday at 9:00 am is:

```
MAILTO="some-admin-email@example.com"
0 9 * * 6 /bin/bash /distillery/backup_all.sh
```

### Backup a single instance

To backup a single instance, you can use the `/distillery/backup_instance.sh` script.
It takes either 1 or 2 arguments:

```bash
# backup a single instance and pick a new file in /backups/final
bash /distillery/backup_instance.sh slug

# backup a single instance into a specific file
bash /distillery/backup_instance.sh /path/to/backup.tar.gz
```

The backup proceeeds as follows:
1. make a copy of the instance configuration
2. shutdown the running instance
3. make a dump of the triplestore and mysql databases
4. make a copy of the file system
5. start the instance again
6. package the data into the final `.tar.gz` file

When uptime is critical, it is possible to skip sets 2 and 5 and leave the instance running.
This might result in inconsistent backup data.
To do so, run the script with the `KEEPALIVE` flag:

```bash
KEEPALIVE=1 bash /distillery/backup_instance.sh slug
```

## SSH Access

The distillery exposes an ssh daemon for users to access individual WissKI Shells.
It is running on port 2222 by default.

To access a shell in a particular barrel set the username equal to the slug.
For instance, to gain access to a shell inside a WissKI instance with a slug `porcelain` use the following command line:

```bash
ssh -p 2222 porcelain@localhost
```

Replace `localhost` with the hostname of the WissKI Distillery.

Inside the container, normal shell acess is provided. 
Both `drush` and `composer` are available. 
No technical reasons using `sudo` or switching to `root` is not possible. 

### Authentication

Authentication is performed using SSH Keys.
Within each instance, ssh keys can be added to the file `/var/www/.ssh/authorized_keys` using the default OpenSSH `authorized_keys` format.

Furthermore, global ssh Keys (that have access to every instance) can be added to a `GLOBAL_AUTHORIZED_KEYS_FILE`. This is set in the Distillery `.env` file, and defaults to `/distillery/authorized_keys/`.

### Port Forwarding

In order to access the __GraphDB Workbench__ or __phpmyadmin__ ssh port forwarding can be used.  
GraphDB is running on the host `triplestore` on port `7200`. 
PhpMyAdmin is running on the host `phpmyadmin` on port `8080`. 

To forward both you can use a command such as:

```bash
ssh -p 2222 -L localhost:7200:triplestore:7200 -L localhost:8080:phpmyadmin:8080 porcelain@localhost
```

This will make GraphDB and PhpMyAdmin available at `localhost:7200` and `localhost:8080` for the duration of the connection. 

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
- Investigate managing phpmyadmin
- Investigate managing graphdb
- Investigate delegating shell access
- Investigate delegating ftp access
- document CNAME structure

<!--  LocalWords:  Vagrantfile vargant phpmyadmin nginx-proxy nginx docker-compose.yml
 -->
<!--  LocalWords:  docker-letsencrypt-nginx-proxy-companion drupal drush Affero graphdb
 -->
<!--  LocalWords:  commerical
 -->
