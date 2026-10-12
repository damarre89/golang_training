#!/bin/bash
# Crosscompile and deploy into RPi2
# go mod init checklist
set -e
TARGET_DEVICE=raspberrypi2.local
WHEREAMI=$(hostname -f)
if [ $WHEREAMI = "raspberrypi3-64" ]; then
    cd /root/projects/golang_training/checklist_server/.
fi
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o checklist .
scp checklist root@$TARGET_DEVICE:/usr/local/bin/checklist.new
ssh root@$TARGET_DEVICE "chmod 0755 /usr/local/bin/checklist.new &&\
    mv /usr/local/bin/checklist.new /usr/local/bin/checklist &&\
    sync &&\
    systemctl restart checklist"