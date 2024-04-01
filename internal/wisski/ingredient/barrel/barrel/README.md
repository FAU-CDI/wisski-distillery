This folder contains all files relevant to a specific WissKI Distillery instance.
This documentation provides basic information about all files and may be updated in future changes.

## Container Files

These are required for the container to function properly.
DO NOT EDIT THESE, changes will be overwritten during the next rebuild.

- `docker-compose.yml`: Defines how the docker image is run
- `Dockerfile`: Dockerfile used to run the actual image
- `.env`: Parameters for `docker-compose.yml`
- `apache.d`, `scripts`, `php.ini.d`: Resources used during the docker-compose build.

To tweak parameters inside the `.env` file, rebuild the image via `wdcli rebuild`.

## Data Files

These contain user data and settings for the instance.
You can edit these (changes will __not__ be overwritten).
Editing these may still risk breaking your instance; do so at your own risk.

- `data`: Drupal Installation and user home directory
- `ssh`: SSH Host keys; modify at your own risk
- `settings.local.php`: settings.php added to this distillery instance
- `prefixes.skip`: If this file exists, prefixes from this wisski instance will not be added to the global resolver.

