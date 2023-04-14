#!/bin/bash

echo +x
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR
cp ../icon-metrics icon-metrics/usr/local/bin/
sed -i "s/Architecture:.*/Architecture: $1/g" icon-metrics/DEBIAN/control
dpkg-deb --root-owner-group --build icon-metrics
