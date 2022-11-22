# WissKI-Distillery

WissKI-Distillery is a Docker-based server provisioning and managing for multiple [WissKI](https://wiss-ki.eu/) instances.

The WissKI Distillery is a set of scripts, tools, and applications that allows to operate a WissKI cloud of distinct but jointly managed WissKI instances, hosted on a dedicated hardware system.
Like the WissKI system, the WissKI Distillery is open source and free to use.

This README contains only technical documentation.
For members of [FAU Erlangen-Nürnberg](https://www.fau.de/) a cloud offering based on this service known as FAUWissKICloud.
Please see https://wisski.data.fau.de/ for related documentation. 

## Overview

This project consists of the following:

- this README containing mostly technical documentation
- a [NEWS.md] file containing recent technical enhancements
- a [Go](https://go.dev/) command `wdcli` 

## High-Level Overview

NOTE: A list of new features can be found in [NEWS](./NEWS.md)

The Distillery consists of a set of instances of WissKIs and (high-level) components.
Components are implemented in the [component](/internal/dis/component) directory.
Furthermore each instance consists of several components refered to as ingredients.
Ingredients are implemented in [ingredient](/internal/wisski/ingredient/) directory.

Each WissKI is implemented as a single [docker](https://www.docker.com/) container that talks to several components.
There are only three components that directly talk to a WissKI Instance
- A [Triplestore](internal/dis/component/triplestore), in this case [GraphDB](http://graphdb.ontotext.com/)
  - a SPARQL backend for WissKI (Version 10.0 or later)
  - Each instance receives a seperate `repository`
- An [SQL Database](internal/dis/component/sql), in this case [MariaDB](https://mariadb.org/)
  - A [phpmyadmin](https://www.phpmyadmin.net/) instance is additionally startyted started on `127.0.0.1:8080`.
  - See [internal/component/sql](internal/dis/component/sql) for implementation details.
- A [SOLR Instance](internal/dis/component/solr)
  - Only preliminary support at this moment

Furthermore to allow end-users to access the WissKIs two further components exist:
- [web](internal/dis/component/web) powered by [Traefik](https://traefik.io/traefik/) to route web traffic to inidividual instances
- A custom [ssh](internal/dis/component/ssh2) server that enables ssh access to individual WissKI instances 

Finally two other components exist:
- A public homepage that lists all instances, and provides basic statistics
- An instance administration system called [info](internal/dis/component/control/info/) that allows web admins to perform certain maintance tasks

# Technical Overview

The go command is almost dependency free. 
It only expects that `docker` and `docker compose` are available.

Each subcommand comes with documentation, which can be found in this readme (and the readme is always outdated), as well as via the command line when passing a `--help` flag.

To bootstrap a new distillery instance, the `wdcli bootstrap` command can be used.
First copy the executable onto the server, using a command similar as:

```bash
GOOS=linux GOARCH=amd64 go build -o wdcli ./cmd/wdcli && scp ./wdcli distillery.example.com:
```

Next, access the server and run the `bootstrap` command:

```
$ ssh distillery.example.com
user@distillery.example.com$ sudo ./wdcli bootstrap
```

This will create a deployment directory (`/var/www/deploy` by default).
Next, edit the configuration file `/var/www/deploy/.env` and customize it to your liking.
Usually it only requires adjustment in very few places.

Next, download a [GraphDB](https://graphdb.ontotext.com/) zip file, and bring the distillery online using:

```bash
/var/www/deploy/wdcli system_update /path/to/graphdb.zip
```
## System Updates

_TLDR: `sudo /var/www/deploy/wdcli system_update /path/to/graphdb.zip`_

To run a WissKI Distillery, several core Docker Instances must be installed.
These are:

- [nginx-proxy](https://github.com/nginx-proxy/nginx-proxy) - an automated nginx reverse proxy

  - This will delegate individual hostnames to appropriate docker containers, see [this blog post](http://jasonwilder.com/blog/2014/03/25/automated-nginx-reverse-proxy-for-docker/) for an overview.
  - Optionally makes use of [docker-letsencrypt-nginx-proxy-companion](https://github.com/nginx-proxy/docker-letsencrypt-nginx-proxy-companion) to automatically provision and renew HTTPS certificates.
  - See [internal/component/web](internal/component/web) for implementation details.

- [MariaDB](https://mariadb.org/) - an SQL server

  - It is configured to run inside a docker container
  - A passwordless `root` account is created, which can only be used from inside the container.
  - An additional admin account (as defined per config file) is created, which is used for administration.
  - A secondary management account is also created. This is configured via the distillery configuration file, and can be access from anywhere.
  - A `bookkeeping` database and table is created by default, to store known WissKI instance metadata in.
  - It is accsssible using `127.0.0.1:3306`
  - A database shell can be opened using `sudo /var/www/deploy/wdcli mysql`.
  - A [phpmyadmin](https://www.phpmyadmin.net/) is started on `127.0.0.1:8080`.
  - See [internal/component/sql](internal/component/sql) for implementation details.


- [proxyssh](https://github.com/tkw1536/proxyssh) - an ssh server that delegates client connections to different WissKIs
  - It is configured to run inside a docker container.
  - Uses a global configurable authorized_keys file.
  - Also allows users to write their own authorized_keys files.
  - See [distillery/resources/compose/ssh](embed/resources/compose/ssh) for implementation details.

- [wdresolve](https://github.com/FAU-CDI/wdresolve) - a global WissKI Distillery Resolver
  - It is configured to run inside a docker container
  - Uses configuration which is updated with `sudo /var/www/deploy/wdcli update_prefix_config` 
  - Running in the browser under the `/go/` path of the main domain.
  - See [distillery/resources/compose/resolver](embed/resources/compose/resolver) for implementation details.

- `dis` - a WissKI Distillery Information Server
  - It is configured to run inside a docker container
  - Running in the browser under the `/dis/` path of the main domain.
  - See [distillery/resources/compose/resolver](embed/resources/compose/dis) for implementation details.

To manage multiple docker containers, this script makes heavy use of [docker compose](https://docs.docker.com/compose/).

Setting up these steps is fully automatic.
In particular, after obtaining a license and the installation zip file for 'GraphDB' one can just run:

```bash
sudo /var/www/deploy/wdcli system_update /path/to/graphdb.zip
```

In principle this script is idempotent, meaning it can be run multiple times achieving the same effect.

## Provisioning a new WissKI instance -- 'wdcli provision'

_TLDR: `sudo /var/www/deploy/wdcli provision name-of-website`_

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
- Sets up a WissKI Salz Adapter to use the GraphDB Repository.

**6. Start the Docker Container**

Finally, we can start the docker container.

These steps can be performed automatically.
To do so, use:

```bash
sudo /var/www/deploy/wdcli provision SLUG
```

## Rebuild an instance -- 'wdcli rebuild'

Sometimes it becomes necessary (because of changes to this project) to rebuild the docker image running a certain docker instance.
To do so, use:

```bash
sudo /var/www/deploy/wdcli rebuild SLUG 
```

Note that rebuilding an instance does restart the docker container resulting in a small (typical < 1 second) interruption to the website in question.
Furthermore, while the container recreated, the old image stays on the host.
To delete all instances, run:

```bash
sudo docker image prune --all
```

To automatically rebuild all instances, use the rebuild command without any arguments:

```bash
sudo /var/www/deploy/wdcli rebuild 
```

## Reserving an instance -- 'wdcli reserve'

Sometimes it is useful to reserve a particular instance name.
This is done by hosting a placeholder website at the domain.
To do so, use:

```bash
sudo /var/www/deploy/wdcli reserve SLUG
```

To un-reserve a website, manually stop the docker stack and remove the folder.

## Purge an existing WissKI instance -- 'wdcli purge'

Sometimes it is required to remove a given WissKI instance.
In particular all parts belonging to it should be removed.

To use it, run:

```bash
sudo /var/www/deploy/wdcli purge SLUG
```

This cannot be undone (expect for manually re-installing a backup or snapshot).
Therefore it typically requires explicit confirmation.

## Open a shell -- 'wdcli shell'

Sometimes manual changes to a given WissKI instance are required.
For this purpose, you can use:

```bash
sudo /var/www/deploy/wdcli shell SLUG
```

This will open a shell in the provided WissKI instance.

## List all instances -- 'wdcli ls'

To list all instances, the following command can be used:

```bash
sudo /var/www/deploy/wdcli ls
```

## Backups & Snapshots -- 'wdcli backup' and 'wdcli snapshot'

### Backup the entire Distillery

This project comes with a backup script.
To make a backup of *all instances*, run:

```bash
sudo /var/www/deploy/wdcli backup
```

Backups may temporarily shutdown individual instances to ensure data consistency.
Typical backup times are a minute or less.

Backups are stored in the `/var/www/deploy/snapshots/archives` directory.
They contain:

- a snapshot of every single instance (see below)
- a complete backup of the SQL database
- nquads of all the GraphDB repositories
- a backup of the configuration + data file(s)

Files are `.tar.gz`ipped.

By default, backups are kept for up to thirty days, after which they are removed.
This can be configured in the WissKI Distillery Configuration File.

### Snapshot a single instance

To snapshot a single instance, you can `sudo /var/www/deploy/wdcli snapshot SLUG`.
It takes either 1 or 2 arguments:

```bash
# snapshot a single instance and pick a new file in /snapshots/archives
sudo /var/www/deploy/wdcli snapshot SLUG

# backup a single instance into a specific file
sudo /var/www/deploy/wdcli snapshot SLUG /path/to/snapshot.tar.gz
```

The snapshot proceeeds as follows:
1. make a copy of the instance configuration
2. shutdown the running instance
3. make a dump of the triplestore and mysql databases
4. make a copy of the file system
5. export all pathbuilders
6. start the instance again
7. package the data into the final `.tar.gz` file

When uptime is critical, it is possible to skip shutting down a running instance.
This might result in inconsistent backup data.
To do so, run the script with the `--keepalive` flag:

```bash
sudo /var/www/deploy/wdcli snapshot SLUG --keepalive
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

### Resolver

In order to resolve WissKI URIs globally, we make use of [wdresolve](https://github.com/FAU-CDI/wdresolve).
This can be queried with a single URI, and will be redirected to the page of the corresponding WissKI Entity.
This is deployed under `/go/` path of the top-level domain.

For example, if the domain name of the distillery instance is `wisski.example.com`, then the resolver would respond to queries like `https://wisski.example.com/go/?uri=https://first.wisski.example.com/content/123`.
The resolver configuration is automatically updated by the `update_prefix_config.sh` script.
It should not be neccessary to reload this configuration manually, as it is automatically called during `system_update.sh`.

It is also possible to manually add a URI prefix to an instance.
For this purpose, add a file named `prefixes` to the base directory of the instance, with one prefix per line.

Furthermore, you can also exclude a specific instance from URL prefix resolving.
This should be the case for cloned or backup instances.
For this purpose, add a file named `prefixes.skip` to the base directory of the instance.
This will casuse the instance to be skipped entirely.

## License

This project and associated files in this repository are licensed as follows:

    WissKI-Distillery - A docker-based WissKI instance server
    Copyright (C) 2020-22 CDI <https://www.cdi.fau.de/>

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
