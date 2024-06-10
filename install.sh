#!/bin/bash

# Serupmon - A simple uptime monitor for your services
# 
# Copyright (c) 2024 Karya Inovasi Anak Bangsa
# Licensed under the MIT License
#
# This script is used to install serupmon on your system.
# It will automatically detect your OS and install the required dependencies.
# 
# This script also can be used to update serupmon to the latest version.
# Or clean up the installation by removing the binary and configuration file. (uninstall)
#
# Usage:
#   bash install.sh [install|update|uninstall]


NC='\033[0m'
RED='\033[0;31m'
BLUE='\033[0;34m'
WHITE='\033[1;37m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BOLD='\033[1m'

ERROR () {
    printf "${RED}${BOLD}%s${NC}\n" "$@"
}

INFO () {
    printf "${BLUE}${BOLD}%s${NC}\n" "$@"
}

WARNING () {
    printf "${YELLOW}${BOLD}%s${NC}\n" "$@"
}

RESULT () {
    printf "${GREEN}${BOLD}=> %s${NC}\n" "$@"
}

STEP () {
    printf "${WHITE}${BOLD}=> %s${NC}\n" "$@"
}

CONFIRM () {
    printf "${YELLOW}${BOLD}"
    read -p "=> $@ [y/N]:" response
    printf "${NC}"
    case $response in
        [yY][eE][sS]|[yY]) 
            true
            ;;
        *)
            false
            ;;
    esac
}

CMD_EXISTS () {
    command -v $1 >/dev/null 2>&1
}

IF_EMPTY_EXIT () {
    if [ -z "$1" ]; then
        ERROR "$2"
        exit 1
    fi

    if [ "$1" == "null" ]; then
        ERROR "$2"
        exit 1
    fi
}

cleanup() {
    echo
    ERROR "Process inetrrupted. Exiting..."
    exit 1
}

# Trap SIGINT
trap cleanup SIGINT

detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$NAME
        VER=$VERSION_ID
    elif type lsb_release >/dev/null 2>&1; then
        OS=$(lsb_release -si)
        VER=$(lsb_release -sr)
    elif [ -f /etc/lsb-release ]; then
        . /etc/lsb-release
        OS=$DISTRIB_ID
        VER=$DISTRIB_RELEASE
    elif [ -f /etc/debian_version ]; then
        OS=Debian
        VER=$(cat /etc/debian_version)
    elif [ -f /etc/SuSe-release ]; then
        OS=SuSE
    elif [ -f /etc/redhat-release ]; then
        OS=RedHat
    else
        OS=$(uname -s)
        VER=$(uname -r)
    fi
}

install_dependencies() {
    if [ "$OS" == "Ubuntu" ] || [ "$OS" == "Debian" ]; then
        sudo apt update -y
        CMD_EXISTS "curl" || sudo apt install -y curl
        CMD_EXISTS "jq" || sudo apt install -y jq
    elif [ "$OS" == "Darwin" ]; then
        # Try detect available package manager
        local package_manager
        if [ -x "$(command -v brew)" ]; then
            package_manager="brew"
        elif [ -x "$(command -v port)" ]; then
            package_manager="port"
        else
            if ! CMD_EXISTS "curl"; then
                ERROR "curl is required to install Homebrew or MacPorts."
                exit 1
            fi

            ERROR "No package manager found. Please install Homebrew or MacPorts first."
            WARNING "Do you want to install Homebrew now?"
            if CONFIRM "Install Homebrew"; then
                /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
                package_manager="brew"
            else
                WARNING "Do you want to install MacPorts now?"
                if CONFIRM "Install MacPorts"; then
                    curl -sSL https://raw.githubusercontent.com/macports/macports-base/master/install.sh | bash
                    package_manager="port"
                else
                    ERROR "Please install Homebrew or MacPorts first."
                    exit 1
                fi
            fi
        fi

        CMD_EXISTS "curl" || $package_manager install curl
        CMD_EXISTS "jq" || $package_manager install jq
    else
        ERROR "Unfortunatelly, your OS is not supported by this script yet."
        STEP "You can install serupmon manually by following the instruction on the README.md"
        exit 1
    fi
}

get_latest_release() {
    local latest_release
    latest_release=$(curl -s https://api.github.com/repos/karyainovasiab/serupmon/releases/latest | jq -r '.tag_name' 2>/dev/null)
    echo $latest_release
}

prepare_serupmon_archive() {
    STEP "Preparing to download the latest release..."
    LATEST_RELEASE=$(get_latest_release)
    IF_EMPTY_EXIT "$LATEST_RELEASE" "Failed to get latest release. Please try again later."
    STEP "Downloading $LATEST_RELEASE for $(uname -s) ($(uname -m))..."

    ARCHIVE_EXT=""
    if [ "$OS" == "Ubuntu" ] || [ "$OS" == "Debian" ]; then
        ARCHIVE_EXT=".deb"
    fi

    ARCHIVE_URL="https://github.com/karyainovasiab/serupmon/releases/download/<TAG>/serupmon_<TAG>_<OS>_<ARCH><EXT>"
    ARCHIVE_URL="${ARCHIVE_URL//<TAG>/$LATEST_RELEASE}"
    ARCHIVE_URL="${ARCHIVE_URL//<OS>/$(uname -s | tr '[:upper:]' '[:lower:]')}"
    ARCHIVE_URL="${ARCHIVE_URL//<ARCH>/$(uname -m | tr '[:upper:]' '[:lower:]')}"
    ARCHIVE_URL="${ARCHIVE_URL//<EXT>/$ARCHIVE_EXT}"

    # Check if the archive is available
    if curl --output /dev/null --silent --head --fail "$ARCHIVE_URL"; then
        curl -L -o serupmon$ARCHIVE_EXT $ARCHIVE_URL
        RESULT "Downloaded serupmon successfully."
    else
        ERROR "Failed to download serupmon. Please try again later."
        exit 1
    fi
}

do_install() {
    if [ "$OS" == "Ubuntu" ] || [ "$OS" == "Debian" ]; then
        sudo dpkg -i serupmon.deb
        # remove the downloaded archive
        rm serupmon.deb
        INFO "Serupmon has been installed successfully."
        sleep 1
        serupmon -h
    elif [ "$OS" == "Darwin" ]; then
        # Detect available path
        local install_path
        local paths=("/usr/local/bin" "/opt/local/bin")
        for path in "${paths[@]}"; do
            if [ -d "$path" ]; then
                install_path=$path
                break
            fi
        done

        if [ -z "$install_path" ]; then
            ERROR "No suitable path found to install serupmon."
            INFO "You can run serupmon by executing the binary directly."
            STEP "TIPS: You can move the binary to a directory in your PATH to run it from anywhere."
            # Check if the binary not executable
            if [ ! -x "serupmon" ]; then
                chmod +x serupmon
            fi

            # Execute the binary
            ./serupmon -h
            exit 0
        fi

        # Check if install_path in PATH
        if ! printf $PATH | grep -q $install_path; then
            WARNING "The install path is not in your PATH. Do you want to add it now?"
            INFO "Please add the following line to your shell configuration file:"
            INFO "export PATH=\$PATH:$install_path"
        fi

        if [ ! -x "serupmon" ]; then
            chmod +x serupmon
        fi

        sudo mv serupmon $install_path
        INFO "Serupmon has been installed successfully."
        sleep 1
        serupmon -h
    else
        ERROR "Unfortunatelly, your OS is not supported by this script yet."
        STEP "You can install serupmon manually by following the instruction on the README.md"
        exit 1
    fi
}

# Install serupmon
install_serupmon() {
    STEP "Starting the installation process..."
    sleep 1
    echo

    STEP "Detecting OS..."
    detect_os
    RESULT "Detected OS: $OS $VER"

    STEP "Installing dependencies..."
    install_dependencies
    RESULT "Dependencies installed successfully."

    prepare_serupmon_archive

    STEP "Installing serupmon..."
    # Ensure the binary is exist
    if [ ! -f "serupmon$ARCHIVE_EXT" ]; then
        ERROR "Failed to download serupmon. Please try again later."
        exit 1
    fi

    do_install
}

# Update serupmon
update_serupmon() {
    STEP "Starting the update process..."
    sleep 1
    echo

    STEP "Detecting OS..."
    detect_os
    RESULT "Detected OS: $OS $VER"

    STEP "Installing dependencies..."
    install_dependencies
    RESULT "Dependencies installed successfully."

    prepare_serupmon_archive

    STEP "Updating serupmon..."
    # Ensure the binary is exist
    if [ ! -f "serupmon$ARCHIVE_EXT" ]; then
        ERROR "Failed to download serupmon. Please try again later."
        exit 1
    fi

    # Backup the old configuration (to be implemented)

    do_install
}

# Uninstall serupmon
uninstall_serupmon() {
    STEP "Starting the uninstallation process..."
    sleep 1
    echo

    STEP "Detecting OS..."
    detect_os
    RESULT "Detected OS: $OS $VER"

    STEP "Uninstalling serupmon..."
    if [ "$OS" == "Ubuntu" ] || [ "$OS" == "Debian" ]; then
        local configs=("/etc/serupmon" "$HOME/.serupmon" "/var/log/serupmon" "/var/run/serupmon")
        for config in "${configs[@]}"; do
            if [ -d "$config" ]; then
                sudo rm -rf $config
            fi
        done

        sudo dpkg -r serupmon
    elif [ "$OS" == "Darwin" ]; then
        # Detect available path
        local install_path
        local paths=("/usr/local/bin" "/opt/local/bin")
        for path in "${paths[@]}"; do
            if [ -d "$path" ]; then
                install_path=$path
                break
            fi
        done

        if [ -z "$install_path" ]; then
            ERROR "No suitable path found to uninstall serupmon."
            exit 1
        fi

        sudo rm $install_path/serupmon
    else
        ERROR "Unfortunatelly, your OS is not supported by this script yet."
        STEP "You can uninstall serupmon manually by following the instruction on the README.md"
        exit 1
    fi

    RESULT "Serupmon has been uninstalled successfully."
}


main() {
    if [ "$EUID" -eq 0 ]; then
        ERROR "Please run this script as a normal user."
        exit 1
    fi

    INFO "Serupmon - A simple uptime monitor for your services"

    case $1 in
        "install")
            install_serupmon
            ;;
        "update")
            update_serupmon
            ;;
        "uninstall")
            uninstall_serupmon
            ;;
        *)
            install_serupmon
            ;;
    esac
}

main "$@"