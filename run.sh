#!/bin/bash

set -e
set -m

if [ -z $1 ]; then
    SLEEP_TIME=3600
else
    SLEEP_TIME=$1
fi

echo "Test if container-events is ready to run ..."
/container-events -test || exit "$?"
echo "Pass the the test"

echo "Starting container event monitor ..."

while [ 1 ]
do
    /container-events &
    sleep ${SLEEP_TIME}
    echo "Restarting container event monitor ..."
    kill %
done
