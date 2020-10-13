#!/bin/bash

current_dir="$(dirname "$0")"

source "${current_dir}/common.sh"

detect_os ()
{
    OS="$(uname -s)"
    ARCH="$(uname -m)"

    case "$OS" in
        Linux)
            if [ -e /etc/lsb-release ]; then
                source /etc/lsb-release

                DIST_ID="${DISTRIB_ID}"
                OS_VERSION="${DISTRIB_RELEASE}"
                OS_CODENAME="${DISTRIB_CODENAME}"
            elif [ -e /etc/os-release]; then
                source /etc/os-release

                DIST_ID="${ID}"
                OS_VERSION="${VERSION_ID}"
                OS_CODENAME="${VERSION_CODENAME}"

            elif [ $(which lsb_release 2>/dev/null) ]; then
                DIST_ID="$(lsb_release -s -i)"
                OS_VERSION="$(lsb_release -s -r)"
                OS_CODENAME="$(lsb_release -s -c)"
            else
                die "Colud not get OS information: there is neither lsb_release tool, nor /etc/lsb-release file"
            fi

            case "$DIST_ID" in
                RedHat*)
                    OS_NAME="RedHat" ;;
                debian)
                    OS_NAME="Debian" ;;
                *)
                    OS_NAME="${DIST_ID}" ;;
            esac
            ;;
        *)
            die "Unsupported OS family: $OS"
            ;;
    esac

    echo "${OS}/${OS_NAME}/${OS_VERSION}/${OS_CODENAME}"
}

detect_os
