#!/bin/sh
# Installation script for deepchain. It tries to move $bin in one of the
# directories stored in $binpaths.

binpaths="/usr/local/bin /usr/bin"
libpaths="/usr/lib /usr/local/lib"


# This variable contains a nonzero length string in case the script fails
# because of missing write permissions.
is_write_perm_missing=""

# Download archive with deepchain binaries according to platform type
PLATFORM=$(uname)
case "$PLATFORM" in
  "Darwin")
    # macOS
    curl -OL  https://github.com/deep-foundaiton/deep-chain/releases/download/v0.2.0/deep_chain_v0.2.0_darwin-amd64.tar.gz
    tar -xzf deep_chain_v0.2.0_darwin-amd64.tar.gz
    for binpath in $binpaths; do
      if cp build_v0.2.0_darwin_amd64/deepchain "$binpath"; then
        for libpath in $libpaths; do
          if cp build_v0.2.0_darwin_amd64/libwasmvm.dylib  "$libpath"; then
            cp build_v0.2.0_darwin_amd64/libwasmvm.so  "$libpath"
            echo "Moved libwasmvm to $libpath"
            break
          else
          if [ -d "$libpath" ] && [ ! -w "$libpath" ]; then
            is_write_perm_missing=1
          fi
          fi
        done
        echo "Moved $bin to $binpath"
        echo "Enjoy your deepchain experience!"
        rm deep_chain_v0.2.0_darwin-amd64.tar.gz
        rm -rf build_v0.2.0_darwin_amd64
        exit 0
      else
      if [ -d "$binpath" ] && [ ! -w "$binpath" ]; then
        is_write_perm_missing=1
        rm deep_chain_v0.2.0_darwin-amd64.tar.gz
        rm -rf build_v0.2.0_darwin_amd64
      fi
      fi
    done
    ;;
      "Linux")
    # Linux distro,
    curl -OL https://github.com/deep-foundaiton/deep-chain/releases/download/v0.2.0/deep-chain_v0.2.0_linux-amd64.tar.gz
    tar -xzf deep_chain_v0.2.0_linux-amd64.tar.gz -C ./
    for binpath in $binpaths; do
      if cp build_v0.2.0_linux_amd64/deepchain "$binpath"; then
        for libpath in $libpaths; do
          if cp build_v0.2.0_linux_amd64/libwasmvm.dylib  "$libpath"; then
            cp build_v0.2.0_linux_amd64/libwasmvm.so  "$libpath"
            echo "Moved libwasmvm to $libpath"
            break
          else
          if [ -d "$libpath" ] && [ ! -w "$libpath" ]; then
            is_write_perm_missing=1
          fi
          fi
        done
        echo "Moved $bin to $binpath"
        echo "Enjoy your deepchain experience!"
        rm deep_chain_v0.2.0_linux-amd64.tar.gz
        rm -rf build_v0.2.0_linux_amd64
        exit 0
      else
      if [ -d "$binpath" ] && [ ! -w "$binpath" ]; then
        is_write_perm_missing=1
        rm deep_chain_v0.2.0_linux-amd64.tar.gz
        rm -rf build_v0.2.0_linux_amd64
      fi
      fi
    done
    ;;
esac


echo "We cannot install $bin in one of the directories $binpaths"

if [ -n "$is_write_perm_missing" ]; then
  echo "It seems that we do not have the necessary write permissions."
  echo "Perhaps try running this script as a privileged user:"
  echo "Or check that you using default library path."
  echo "    sudo $0"
  echo
fi

exit 1
