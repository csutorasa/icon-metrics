#!/usr/bin/env sh

STATUS=$(wget -q -O - http://127.0.0.1:8080/status)
if [ "$?" != "0" ]; then
    exit 1
fi
if [ "$STATUS" != "OK" ]; then
    exit 1
fi
