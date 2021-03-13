#!/bin/bash

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

    echo "Changing SUID for agent bin file ..."

    if chmod u+s "$PWD/r2dtools"; then
        echo "Agent bin SUID is successfully changed"
    else
        die "Could not change SUID for agent bin file."
    fi
}
