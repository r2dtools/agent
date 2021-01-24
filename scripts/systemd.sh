#!/bin/bash

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
