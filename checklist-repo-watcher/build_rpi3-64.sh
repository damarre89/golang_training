#!/bin/bash
# build from amd64 into aarch64 target
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o checklist-repo-watcher .
#cp checklist-repo-watcher /usr/local/bin/.
#chmod 0755 /usr/local/bin/checklist-repo-watcher
#sync