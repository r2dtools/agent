#!/bin/bash

CURRENT_DIR="$(dirname "$0")"

source "${CURRENT_DIR}/common.sh"

# remove r2dtools agent from systemd
remove_systemd_service()
{
    echo "Disabling R2DTools agent systemd service ..."

    if sudo systemctl disable "r2dtools"; then
        echo "R2DTools agent systemd service is disabled."
    else
        die "Could not disable R2DTools agent systemd service."
    fi

    sudo rm "/etc/systemd/system/r2dtools.service" &> /dev/null
    sudo systemctl daemon-reload &> /dev/null
}

# remove r2dtools:r2dtools agent group
remove_agent_group()
{
    sudo userdel $USER
    sudo groupdel $GROUP
}

source "${CURRENT_DIR}/systemd.sh"
stop_systemd_service
remove_systemd_service
remove_agent_group

echo "R2DTools agent was successfully uninstalled."
echo "Now you can remove agent directory."
