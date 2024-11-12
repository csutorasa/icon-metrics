#!/bin/bash

ICON_METRICS_VERSION="1.3.1"
ICON_METRICS_ARCH=${1:-amd64}

ICON_METRICS_INSTALL_TMP_DIR=~/.icon-metrics-install
ICON_METRICS_DOWNLOAD_FILENAME=icon-metrics-linux-$ICON_METRICS_ARCH
ICON_METRICS_DOWNLOAD_URL="https://github.com/csutorasa/icon-metrics/releases/download/$ICON_METRICS_VERSION/$ICON_METRICS_DOWNLOAD_FILENAME.zip"

if [[ -d "$ICON_METRICS_INSTALL_TMP_DIR" ]] ; then
    echo "Temp directory $ICON_METRICS_INSTALL_TMP_DIR exists!"
    exit 1
fi

if which curl >/dev/null ; then
    mkdir -p $ICON_METRICS_INSTALL_TMP_DIR
    curl -L $ICON_METRICS_DOWNLOAD_URL --output "$ICON_METRICS_INSTALL_TMP_DIR/$ICON_METRICS_DOWNLOAD_FILENAME.zip"
elif which wget >/dev/null ; then
    mkdir -p $ICON_METRICS_INSTALL_TMP_DIR
    wget $ICON_METRICS_DOWNLOAD_URL -P "$ICON_METRICS_INSTALL_TMP_DIR"
else
    echo "Cannot download zip, neither curl nor wget command is available!"
    exit 1
fi
if [ $? -ne 0 ] ; then
    echo "Failed to download file!"
    exit 1
fi

if which 7z >/dev/null ; then
    7z e "$ICON_METRICS_INSTALL_TMP_DIR/$ICON_METRICS_DOWNLOAD_FILENAME.zip" -o"$ICON_METRICS_INSTALL_TMP_DIR"
elif which unzip >/dev/null ; then
    unzip "$ICON_METRICS_INSTALL_TMP_DIR/$ICON_METRICS_DOWNLOAD_FILENAME" -d "$ICON_METRICS_INSTALL_TMP_DIR"
else
    echo "Cannot unzip, neither 7z nor unzip command is available!"
    exit 1
fi
if [ $? -ne 0 ] ; then
    echo "Failed to unzip file!"
    rm -rf "$ICON_METRICS_INSTALL_TMP_DIR"
    exit 1
fi
mkdir -p /usr/local/bin/
mv "$ICON_METRICS_INSTALL_TMP_DIR/icon-metrics" /usr/local/bin/
chmod +x /usr/local/bin/icon-metrics
mkdir -p /etc/icon-metrics
ICON_METRICS_CONFIG_MESSAGE="Check your existing config"
if [[ ! -f /etc/icon-metrics/config.yml ]] ; then
    ICON_METRICS_CONFIG_MESSAGE="Add your device to the example config"
    mv "$ICON_METRICS_INSTALL_TMP_DIR/config.yml" /etc/icon-metrics/
fi
mkdir -p /var/log/icon-metrics
rm -rf "$ICON_METRICS_INSTALL_TMP_DIR"
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
    echo "$ICON_METRICS_CONFIG_MESSAGE"
    echo -e "\t/etc/icon-metrics/config.yml"
    echo "Enable and start the service"
    echo -e "\tsystemctl daemon-reload"
    echo -e "\tsystemctl enable --now icon-metrics.service"
else
    echo "$ICON_METRICS_CONFIG_MESSAGE"
    echo -e "\t/etc/icon-metrics/config.yml"
    echo "Start the application"
    echo -e "\t/usr/local/bin/icon-metrics --config /etc/icon-metrics/config.yml > /var/log/icon-metrics/log.log"
fi
