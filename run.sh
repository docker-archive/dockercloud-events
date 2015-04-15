#!/bin/bash

set -e
set -m

export DOCKER_BINARY="/docker"

eval "${DOCKER_BINARY} version" >/dev/null 2>&1  || {
	echo "Cannot run docker at ${DOCKER_BINARY}" ;
	exit 1;
}

if [ "${TUTUM_AUTH}" == "**None**" ]; then
	echo "Need to specify TUTUM_AUTH"
    exit 1
fi

if [ "${NODE_UUID}" == "**None**" ]; then
	echo "Need to specify NODE_UUID"
    exit 1
fi

echo "Testing execution environment"
/container-events -test

echo "Starting container event monitor ..."
while [ 1 ]
do
    /container-events &
    sleep ${RESTART_INTERVAL}
    echo "Restarting container event monitor ..."
    kill %
done
