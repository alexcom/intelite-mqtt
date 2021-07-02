#!/usr/bin/env bash
rm ./intelite-mqtt || true
CGO_ENABLED=false GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags="-static"' -o intelite-mqtt ./cmd/main.go
