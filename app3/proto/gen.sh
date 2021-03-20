#!/bin/bash
set -o nounset
set -o pipefail
# set -x

DIR="$( cd "$( dirname "$0" )" && pwd -P )"
cd $DIR


protoc -I ./ --go_out=plugins=grpc:../pkg/helloworld ./helloworld.proto