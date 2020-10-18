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

