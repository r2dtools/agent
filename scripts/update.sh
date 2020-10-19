#!/bin/bash

CURRENT_DIR="$(dirname "$0")"

source "${CURRENT_DIR}/common.sh"

URL="http://r2dtools.com/builds/r2dtools-latest.tar.gz"

# start r2dtools agent service
start_systemd_service()
{
    echo "Starting R2DTools agent service ..."
    systemctl start "r2dtools"
    
    if systemctl status "r2dtools"; then
        systemctl enable "r2dtools"
        echo "R2DTools agent service successfully started."
    else
        die "Could not start R2DTools agent service."
    fi
}

# download and unpack agent
download_and_unpack_agent()
{
    local unpackedDirName="r2dtools-agent"
    local filename="r2dtools-agent.tar.gz"
    local directory="/tmp"
    local filePath="${directory}/${filename}"
    local unpackedDirPath="${directory}/${unpackedDirName}"
    local PWD=$(pwd)
    
    rm $filePath &> /dev/null
    rm -r $unpackedDirPath &> /dev/null

    echo "Downloading the latest version of R2DTools agent ..."
    
    if wget -O $filePath $URL; then
        echo "R2DTools agent is sucessfully downloaded."
    else
        die "Could not download R2DTools agent."
    fi
    
    echo "Creating directory ${unpackedDirPath} ..."

    if mkdir "${unpackedDirPath}"; then
        echo "Directory ${unpackedDirPath} is successfully created."
    else
        die "Could not create directory ${unpackedDirPath}."
    fi

    echo "Unpacking R2DTools agent ..."

    if tar -xzf ${filePath} -C "${unpackedDirPath}"; then
        echo "R2DTools agent is successfully unpacked."
    else
        die "Could not unpack R2DTools agent."
    fi

    echo "Copying files to ${PWD} ..."

    if cp -a "${unpackedDirPath}/." ${PWD}; then
        echo "R2DTools agent files successfully copied."
    else
        die "Could not copy R2DTools agent files to ${PWD}."
    fi
}

# update r2dtools agent
update()
{
    stop_systemd_service
    download_and_unpack_agent
    set_agent_dir_owner
    start_systemd_service
}

update
