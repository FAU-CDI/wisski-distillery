#!/bin/bash
set -e

# This file will load all the library functions needed by the various scripts. 
# It should be 'source'd only, if it is not we bail out here. 
if [[ "$0" = "$BASH_SOURCE" ]]; then
   echo "This file should not be executed directly, it should be 'source'd only. "
   exit 1;
fi

# Set variables for the script_dir and the lib_dir
SCRIPT_DIR="$(pwd)"
LIB_DIR="$SCRIPT_DIR/lib"

# Next, we load a bunch of utility functions stored in lib/lib_<number>_<system>.sh 
# These contain functionality used in the various scripts. 
source "$LIB_DIR/00_init.sh";
source "$LIB_DIR/10_config.sh";
source "$LIB_DIR/20_slug.sh";
source "$LIB_DIR/30_utils.sh";
