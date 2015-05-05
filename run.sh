#!/bin/sh

set -e

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

echo "Starting container event monitor ..."
exec /container-events 
