#!/bin/bash

CURRENT_DIR="$(dirname "$0")"

source "${CURRENT_DIR}/common.sh"

URL="http://r2dtools.com/builds/r2dtools-latest.tar.gz"

# copy agent files
copy_agent_files()
{
    local PWD=$(pwd)

    echo "Copying files to ${TARGET_DIR} ..."

    if cp -p r2dtools lego .version ${TARGET_DIR}; then
        echo "R2DTools agent files successfully copied."
    else
        die "Could not copy R2DTools agent files to ${TARGET_DIR}."
    fi

    if cp -p -R scripts ${TARGET_DIR}; then
        echo "R2DTools agent scritps successfully copied."
    else
        die "Could not copy R2DTools agent scripts to ${TARGET_DIR}."
    fi
}

source "${CURRENT_DIR}/systemd.sh"
stop_systemd_service
copy_agent_files
source "${CURRENT_DIR}/permissions.sh"
set_agent_dir_owner
source "${CURRENT_DIR}/post_update.sh"
start_systemd_service
