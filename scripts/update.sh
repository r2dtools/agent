#!/bin/bash

CURRENT_DIR="$(dirname "$0")"

source "${CURRENT_DIR}/common.sh"

URL="http://r2dtools.com/builds/r2dtools-latest.tar.gz"

# copy agent files
copy_agent_files()
{
    local PWD=$(pwd)

    # do not replace user`s config file.
    rm "${PWD}/config/params.yaml" &> /dev/null

    echo "Copying files to ${TARGET_DIR} ..."

    if cp -a "${PWD}/." ${TARGET_DIR}; then
        echo "R2DTools agent files successfully copied."
    else
        die "Could not copy R2DTools agent files to ${TARGET_DIR}."
    fi
}

source "${CURRENT_DIR}/systemd.sh"
stop_systemd_service
copy_agent_files
source "${CURRENT_DIR}/permissions.sh"
set_agent_dir_owner
source "${CURRENT_DIR}/post_update.sh"
start_systemd_service
