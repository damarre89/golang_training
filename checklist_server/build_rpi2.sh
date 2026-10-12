#!/bin/bash
# Crosscompile for RPi2 arch
# go mod init checklist
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o checklist .