#!/bin/bash

read -p "version: " version
if [ -z "$version" ]; then
    echo "version is required."
    exit 1
fi

export GOOS=windows
export GOARCH=amd64
go build -ldflags "-H=windowsgui" -o "build/clock-out-reminder.$version.exe"

export GOARCH=arm64
go build -ldflags "-H=windowsgui" -o "build/clock-out-reminder.$version.arm64.exe"
