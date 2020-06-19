#!/bin/bash

current_dir="$(dirname "$0")"

source "${current_dir}/common.sh"
source "${current_dir}/os.sh"

# Check that the current platform is supported
check_arch()
{
    echo "Checking platform type ..."
    ARCH="$(uname -m)"
    echo "${ARCH}"

    if [ "$ARCH" != "x86_64" ]; then
        die "Unsupported platform: ${ARCH}"
    fi

    echo "OK. The current platform is supported."
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

    echo "OK. The current OS is supported."

}

install()
{
    check_arch
    check_os
}

install

