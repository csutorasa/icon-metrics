#!/bin/bash

ICON_METRICS_VERSION="1.0.1"
ICON_METRICS_ARCH=${1:-amd64}

INSTALL_TMP_DIR=~/.icon-metrics-install
DOWNLOAD_URL="https://github.com/csutorasa/icon-metrics/releases/download/$ICON_METRICS_VERSION/icon-metrics-linux-$ICON_METRICS_ARCH.zip"
mkdir -p $INSTALL_TMP_DIR
if which curl >/dev/null ; then
    curl $DOWNLOAD_URL --output $INSTALL_TMP_DIR
elif which wget >/dev/null ; then
    wget $DOWNLOAD_URL -P $INSTALL_TMP_DIR
else
    echo "Cannot download zip, neither curl nor wget command is available"
    exit 1
fi
if ! which unzip >/dev/null ; then 
    echo "Cannot unzip, unzip command is available"
    exit 1
fi
unzip $INSTALL_TMP_DIR/icon-metrics-*.zip -d $INSTALL_TMP_DIR
mkdir -p /usr/local/bin/
mv $INSTALL_TMP_DIR/icon-metrics /usr/local/bin/
chmod +x /usr/local/bin/icon-metrics
mkdir -p /etc/icon-metrics
mv $INSTALL_TMP_DIR/config.yml /etc/icon-metrics/
rm -rf $INSTALL_TMP_DIR
if which systemctl >/dev/null ; then
    car /etc/systemd/system/icon-metrics.service << EOF
[Unit]
Description=iCON metrics publisher

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
ExecStart=/usr/local/bin/icon-metrics --config /etc/icon-metrics/config.yml
WorkingDirectory=/usr/local/bin/
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
    echo -e "\t/usr/local/bin/icon-metrics --config /etc/icon-metrics/config.yml"
fi
