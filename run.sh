#!/bin/bash

set -e
set -m

DOCKER_BINARY=/docker

eval "${DOCKER_BINARY} version" >/dev/null 2>&1  ||( echo "Cannot run docker binary at {DOCKER_BINARY}" && exit 1)

if [ -z $1 ]; then
    SLEEP_TIME=3600
else
    SLEEP_TIME=$1
fi

echo "Test if container-events is ready to run ..."
/container-events -test -dockerBinary=${DOCKER_BINARY} || exit "$?"
echo "Pass the the test"

echo "Starting container event monitor ..."

while [ 1 ]
do
    /container-events -dockerBinary=${DOCKER_BINARY} &
    sleep ${SLEEP_TIME}
    echo "Restarting container event monitor ..."
    kill %
done
