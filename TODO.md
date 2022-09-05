# WissKI-Distillery in Go

This document describes the process of moving the distillery to using golang for the control plane (outside of docker containers).

## Bootstrapping

This documents the bootstraping process.
Work in progress.

- `wdcli bootstrap $DIRECTORY`
    0. Create the deployment directory
    1. Copy over the executable (unless it already exists)
    2. Create a default configuration file (unless it already exists)
    3. Store the directory in a file called .wdcli in the $HOME directory

- `wdcli system_update`
    - to be documented
## Future Work

- Move `provision_entrypoint.sh` into go
- Rename backups to 'snapshots' and make them restorable
    - Snapshot the docker images being used also!
- Avoid running `docker compose` executable and shift it to a library
- Automatically bootstrap the docker container sql connection (use proper environment variables)
- Make error handling consistent
- Add a server that serves information
- Migrate the individual commands below
- restructure resource files
- Documentation

## Migrating Individual Commands
- [ ] backup_all.sh
- [x] backup_instance.sh
- [x] blind_update.sh
- [x] blind_update_all.sh
- [x] cron_all.sh
- [x] info.sh
- [x] ls.sh
- [x] make_mysql_account.sh
- [ ] monday_full.sh
- [ ] monday_short.sh
- [x] mysql.sh
- [x] provision.sh
- [x] purge.sh
- [x] rebuild.sh
- [x] rebuild_all.sh
- [x] reserve.sh
- [x] shell.sh
- [x] system_install.sh
- [x] system_update.sh
- [x] update_prefix_config.sh

## TO BE REMOVED
- [ ] call_update_php_hack.sh
