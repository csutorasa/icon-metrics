#!/bin/bash

cp ../icon-metrics icon-metrics/usr/local/bin/
sed -i "s/Architecture:.*/Architecture: $1/g"
dpkg-deb --root-owner-group --build icon-metrics
