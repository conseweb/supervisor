#!/bin/bash

set -ex

ROOT="$( cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd )"

cd $ROOT/integration-tests
go test -check.v github.com/conseweb/supervisor/integration-tests