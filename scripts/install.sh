#!/bin/bash

CURRENT_DIR="$(dirname "$0")"

source "${CURRENT_DIR}/common.sh"
source "${CURRENT_DIR}/os.sh"

USER="r2dtools"
GROUP="r2dtools"

# Check that the current platform is supported
check_arch()
{
    echo "Checking platform type ..."
    ARCH="$(uname -m)"
    echo "${ARCH}"

    if [ "$ARCH" != "x86_64" ]; then
        die "Unsupported platform: ${ARCH}"
    fi

    echo "OK. The current platform ${ARCH} is supported."
}

# Check that the current OS is supported
check_os()
{
    DEBIAN="Debian"
    UBUNTU="Ubuntu"
    CENTOS="CentOS"

    echo "Detecting OS type and version ..."
    detect_os

    case "$OS_NAME" in
        "$DEBIAN") ;;
        "$UBUNTU") ;;
        "$CENTOS") ;;
        *)
            die "Unsupported OS: ${OS_NAME}." ;;
    esac

    echo "OK. The current OS ${OS_NAME} is supported."

}

# creates user/group r2dtools/r2dtools if it does not exist yet
create_user_group()
{
    if grep -q $GROUP "/etc/group"; then
        echo "Group '${GROUP}' already exists."
    else
        if groupadd $GROUP; then
            echo "Group '${GROUP}' successfully created."
        else
            die "Could not create group '${GROUP}'."
        fi
    fi

    if id $USER &> /dev/null; then
        echo "User '${USER}' is already exists."
    else
        if useradd -g $GROUP $USER; then
            echo "User '${USER}' successfully created."
        else
            die "Could not create user '${USER}'."
        fi
    fi
}

# set correct owner for agent directory
set_agent_dir_owner()
{
    local PWD=$(pwd)
    
    echo "Set owner $USER:$GROUP for agent directory ..."
    
    if chown -R $USER:$GROUP $PWD; then
        echo "Agent directory owner is successfully changed to $USER:$GROUP"
    else
        die "Could not change $USER:$GROUP owner for agent directory."
    fi
}

# create r2dtools agent systemd service
create_systemd_service()
{
    local SERVICE_FILE="/etc/systemd/system/r2dtools.service"
    local PWD=$(pwd)
    cp "${CURRENT_DIR}/r2dtools.service" ${SERVICE_FILE}
    sed -i "s/R2DTOOLS_USER/${USER}/" ${SERVICE_FILE}
    sed -i "s#R2DTOOLS_SERVE#${PWD}/r2dtools serve#" ${SERVICE_FILE}
    
    if systemctl start "r2dtools"; then
        echo "R2DTools agent service successfully started."
    else
        die "Could not start R2DTools agent service."
    fi

    systemctl enable "r2dtools"
}

install()
{
    check_arch
    check_os
    create_user_group
    set_agent_dir_owner
    create_systemd_service
}

install
