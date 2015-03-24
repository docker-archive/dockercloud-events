#!/bin/bash

set -e
set -m

DOCKER_BINARY=/docker

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

/container-events -test \
	-dockerBinary "${DOCKER_BINARY}" \
	-dockerHost "${DOCKER_HOST}" \
	-tutumHost "${TUTUM_HOST}" \
	-auth "${TUTUM_AUTH}" \
	-uuid "${NODE_UUID}" || \
	exit "$?"

echo "Starting container event monitor ..."
while [ 1 ]
do
    /container-events \
		-dockerBinary "${DOCKER_BINARY}" \
		-dockerHost "${DOCKER_HOST}" \
		-tutumHost "${TUTUM_HOST}" \
		-auth "${TUTUM_AUTH}" \
		-uuid "${NODE_UUID}" &
    sleep ${RESTART_INTERVAL}
    echo "Restarting container event monitor ..."
    kill %
done
