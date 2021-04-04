#!/bin/bash

USER="r2dtools"
GROUP="r2dtools"
TARGET_DIR="/opt/r2dtools"

die()
{
    echo "ERROR: $*" >&2
	exit 1
}
