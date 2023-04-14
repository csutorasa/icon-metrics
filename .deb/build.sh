#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR
cp ../icon-metrics icon-metrics/usr/local/bin/icon-metrics
sed -i "s/Architecture:.*/Architecture: $1/g" icon-metrics/DEBIAN/control
dpkg-deb --root-owner-group --build icon-metrics
