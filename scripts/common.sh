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

    echo "Changing owner root:${GROUP} for agent bin file ..."

    if chown "root:${GROUP}" "$PWD/r2dtools"; then
        echo "Agent bin owner is successfully changed to root:${GROUP}"
    else
        die "Could not change root:$GROUP owner for agent bin file."
    fi

    echo "Changing SUID for agent bin file ..."

    if chmod u+s "$PWD/r2dtools"; then
        echo "Agent bin SUID is successfully changed"
    else
        die "Could not change SUID for agent bin file."
    fi
}

die()
{
    echo "ERROR: $*" >&2
	exit 1
}
