#!/bin/bash

USER="r2dtools"
GROUP="r2dtools"

# set correct owner for agent directory
set_agent_dir_owner()
{
    local PWD=$(pwd)
    
    echo "Changing owner $USER:$GROUP for agent directory ..."
    
    if chown -R $USER:$GROUP $PWD; then
        echo "Agent directory owner is successfully changed to $USER:$GROUP"
    else
        die "Could not change $USER:$GROUP owner for agent directory."
    fi
}

die()
{
    echo "ERROR: $*" >&2
	exit 1
}

# stop r2dtools agent service
stop_systemd_service()
{
    echo "Stoping R2DTools agent service ..."

    if systemctl stop "r2dtools"; then
        echo "R2DTools agent service is successfully stoped."
    else
        die "Could not stop R2DTools agent service."
    fi
}
