#!/bin/sh
# Installation script for space-pussy.

binpaths="/usr/local/bin /usr/bin"
libpaths="/usr/lib /usr/local/lib"


# This variable contains a nonzero length string in case the script fails
# because of missing write permissions.
is_write_perm_missing=""

# Download and install space-pussy from GitHub
PLATFORM=$(uname)
case "$PLATFORM" in
  "Darwin"|"Linux")
    wget https://github.com/greatweb/space-pussy/archive/refs/tags/v0.0.3.zip
    unzip v0.0.3.zip
    cd space-pussy-0.0.3/
    make build
    for binpath in $binpaths; do
      if cp build/pussy "$binpath"; then
        echo "Moved pussy to $binpath"
        echo "Enjoy your space-pussy experience!"
        rm ~/v0.0.3.zip
        rm -rf ~/space-pussy-0.0.3
        exit 0
      else
        if [ -d "$binpath" ] && [ ! -w "$binpath" ]; then
          is_write_perm_missing=1
          rm ~/v0.0.3.zip
          rm -rf ~/space-pussy-0.0.3
        fi
      fi
    done
    ;;
esac

echo "We cannot install pussy in one of the directories $binpaths"

if [ -n "$is_write_perm_missing" ]; then
  echo "It seems that we do not have the necessary write permissions."
  echo "Perhaps try running this script as a privileged user:"
  echo "Or check that you are using the default library path."
  echo "    sudo $0"
  echo
fi

exit 1
