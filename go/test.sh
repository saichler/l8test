#!/usr/bin/env bash

# Fail on errors and don't open cover file
set -e

# clean up
git checkout go.mod
rm -rf go.sum
rm -rf vendor

# fetch dependencies
GOPROXY=direct GOPRIVATE=github.com go mod tidy

./build-security.sh

go mod vendor

# Run unit tests with coverage
go test -tags=unit -v -coverpkg=./infra/... -coverprofile=cover.html ./... --failfast

# Open the coverage report in a browser
go tool cover -html=cover.html
