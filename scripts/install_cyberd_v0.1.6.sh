#!/bin/sh
# Installation script for cyber. It tries to move $bin in one of the
# directories stored in $binpaths.

INSTALL_DIR=$(dirname $0)

binpaths="/usr/local/bin /usr/bin"


# This variable contains a nonzero length string in case the script fails
# because of missing write permissions.
is_write_perm_missing=""

# Download archive with cyberdcli binaries according to platform type
PLATFORM=$(uname)
case "$PLATFORM" in
  "Darwin")
    # macOS
    wget https://mars.cybernode.ai/go-cyber/go_cyber_v0.1.6_darwin-amd64.tar.gz
    for binpath in $binpaths; do
      if tar -xzf go_cyber_v0.1.6_darwin-amd64.tar.gz -C "$binpath"; then
        echo "Moved $bin to $binpath"
        rm go_cyber_v0.1.6_darwin-amd64.tar.gz
        exit 0
      else
      if [ -d "$binpath" ] && [ ! -w "$binpath" ]; then
        is_write_perm_missing=1
        rm go_cyber_v0.1.6_darwin-amd64.tar.gz
      fi
      fi
    done
    ;;
  "Linux")
    # Linux distro,
    wget https://mars.cybernode.ai/go-cyber/go_cyber_v0.1.6_linux-amd64.tar.gz
    for binpath in $binpaths; do
      if tar -xzf go_cyber_v0.1.6_linux-amd64.tar.gz -C "$binpath"; then
        echo "Moved $bin to $binpath"
        rm go_cyber_v0.1.6_linux-amd64.tar.gz
        exit 0
      else
      if [ -d "$binpath" ] && [ ! -w "$binpath" ]; then
        is_write_perm_missing=1
        rm go_cyber_v0.1.6_darwin-amd64.tar.gz
      fi
      fi
    done
    ;;
esac


echo "We cannot install $bin in one of the directories $binpaths"

if [ -n "$is_write_perm_missing" ]; then
  echo "It seems that we do not have the necessary write permissions."
  echo "Perhaps try running this script as a privileged user:"
  echo
  echo "    sudo $0"
  echo
fi

exit 1
