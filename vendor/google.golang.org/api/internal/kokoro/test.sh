#!/bin/bash

# Copyright 2018 Google LLC.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Fail on any error
set -eo pipefail

export GOOGLE_APPLICATION_CREDENTIALS="${KOKORO_GFILE_DIR}/secret_manager/go-cloud-integration-service-account"

# Display commands being run
set -x

# cd to project dir on Kokoro instance
cd github/google-api-go-client

go version

# Set $GOPATH
export GOPATH="$HOME/go"
export GOCLOUD_HOME=$GOPATH/src/google.golang.org/api/
export PATH="$GOPATH/bin:$PATH"
export GO111MODULE=on
mkdir -p $GOCLOUD_HOME

# Move code into $GOPATH and get dependencies
git clone . $GOCLOUD_HOME
cd $GOCLOUD_HOME

try3() { eval "$*" || eval "$*" || eval "$*"; }

# All packages, including +build tools, are fetched.
try3 go mod download
./internal/kokoro/vet.sh

# Testing the generator itself depends on a generation step
cd google-api-go-generator
go generate
cd ..

# Run tests and tee output to log file, to be pushed to GCS as artifact.
if [[ $KOKORO_JOB_NAME == *"continuous"* ]]; then
    go test -race -v ./... 2>&1 | tee $KOKORO_ARTIFACTS_DIR/$KOKORO_GERRIT_CHANGE_NUMBER.txt
else
    go test -race -v -short ./... 2>&1 | tee $KOKORO_ARTIFACTS_DIR/$KOKORO_GERRIT_CHANGE_NUMBER.txt
fi
