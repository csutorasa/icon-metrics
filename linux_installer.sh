#!/bin/bash

ICON_METRICS_VERSION="1.1.0"
ICON_METRICS_ARCH=${1:-amd64}

INSTALL_TMP_DIR=~/.icon-metrics-install
DOWNLOAD_FILENAME=icon-metrics-linux-$ICON_METRICS_ARCH.zip
DOWNLOAD_URL="https://github.com/csutorasa/icon-metrics/releases/download/$ICON_METRICS_VERSION/$DOWNLOAD_FILENAME"

if [[ -d "$INSTALL_TMP_DIR" ]] ; then
    echo "Temp directory $INSTALL_TMP_DIR exists!"
    exit 1
fi
if ! which unzip >/dev/null ; then 
    echo "unzip command is available!"
    exit 1
fi

if which curl >/dev/null ; then
    mkdir -p $INSTALL_TMP_DIR
    curl $DOWNLOAD_URL --output "$INSTALL_TMP_DIR/$DOWNLOAD_FILENAME"
elif which wget >/dev/null ; then
    mkdir -p $INSTALL_TMP_DIR
    wget $DOWNLOAD_URL -P "$INSTALL_TMP_DIR"
else
    echo "Cannot download zip, neither curl nor wget command is available!"
    exit 1
fi
if [ $? -ne 0 ] ; then
    echo "Failed to download file!"
    exit 1
fi
unzip "$INSTALL_TMP_DIR/$DOWNLOAD_FILENAME" -d "$INSTALL_TMP_DIR"
if [ $? -ne 0 ] ; then
    echo "Failed to unzip file!"
    rm -rf "$INSTALL_TMP_DIR"
    exit 1
fi
mkdir -p /usr/local/bin/
mv "$INSTALL_TMP_DIR/icon-metrics" /usr/local/bin/
chmod +x /usr/local/bin/icon-metrics
mkdir -p /etc/icon-metrics
if [[ ! -f /etc/icon-metrics/config.yml ]] ; then
    mv "$INSTALL_TMP_DIR/config.yml" /etc/icon-metrics/
fi
mkdir -p /var/log/icon-metrics
rm -rf "$INSTALL_TMP_DIR"
if which systemctl >/dev/null ; then
    cat /etc/systemd/system/icon-metrics.service << EOF
[Unit]
Description=iCON metrics publisher

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
ExecStart=/usr/local/bin/icon-metrics --config /etc/icon-metrics/config.yml
WorkingDirectory=/usr/local/bin/
StandardOutput=/var/log/icon-metrics/log.log
StandardError=/var/log/icon-metrics/error.log
Restart=always
EOF
    echo "Edit the config"
    echo -e "\t/etc/icon-metrics/config.yml"
    echo "Enable and start the service"
    echo -e "\tsystemctl daemon-reload && systemctl enable --now icon-metrics.service"
else
    echo "Edit the config"
    echo -e "\t/etc/icon-metrics/config.yml"
    echo "Start the application"
    echo -e "\t/usr/local/bin/icon-metrics --config /etc/icon-metrics/config.yml > /var/log/icon-metrics/log.log"
fi
