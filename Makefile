# SPDX-FileCopyrightText: 2021 Eric Neidhardt
# SPDX-License-Identifier: CC0-1.0

all: test build

build:
	go build ./cmd/webserver/

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o out/ ./cmd/webserver/

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o out/ ./cmd/webserver/

test:
	go test ./...
