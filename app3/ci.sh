#!/bin/bash
set -o nounset
set -o pipefail
# set -x

DIR="$( cd "$( dirname "$0" )" && pwd -P )"
cd $DIR


DOCKER_REPO="rival0605"
image_name="jaeger"
image_tag="app3"



docker build -f ./Dockerfile -t $DOCKER_REPO/$image_name:$image_tag ./ || exit

docker push $DOCKER_REPO/$image_name:$image_tag
