#!/bin/bash

echo +x
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR
gzip --best -n icon-metrics/usr/share/man/man1/icon-metrics.1
gzip --best -n icon-metrics/usr/share/doc/icon-metrics/changelog.Debian
