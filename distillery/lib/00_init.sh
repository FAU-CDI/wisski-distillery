#!/bin/bash
set -e

# This is a library file. 
# It should be 'source'd only, if it is not we bail out here. 
if [[ "$0" = "$BASH_SOURCE" ]]; then
   echo "This file should not be executed directly, it should be 'source'd only. "
   exit 1;
fi

# Check that we are running on a linux system, to prevent accidentally running on Windows or Mac. 
# Ideally we would want to explicitly limit to Debian / Ubuntu but at this point that's not needed. 
OS="$(uname -s)"
case "${OS}" in
    Linux*) :;;
    *)      echo "This script must be run under Linux. "; exit 1;;
esac

# To prevent accidentally messing up permissions, we need to always run as root. 
# Check that the uid is 0, and otherwise bail out. 
if [[ $EUID -ne 0 ]]; then
   echo "This script should be run as root, use 'sudo' if in doubt. "
   exit 1;
fi

# We enable shell aliases, to allow us to setup utility functions much more easily. 
# To be safe that we don't have any other ones in the environment, we first unalias everything. 
unalias -a
shopt -s expand_aliases

# Setup some basic input/output functions
function log_info() {
   echo -e "\033[1m$1\033[0m"
}

function log_ok() {
   echo -e "\033[0;32m$1\033[0m"
}

function log_warn() {
   echo -e "\033[1;33m$1\033[0m"
}

function log_error() {
   echo -e "\033[0;31m$1\033[0m"
}

if [ -n "$DISABLE_LOG" ]; then
   function log_info() {
      true
   }
   function log_ok() {
      true
   }
   function log_warn() {
      true
   }
   function log_error() {
      true
   }
fi