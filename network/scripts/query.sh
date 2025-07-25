#!/bin/bash

if [ $# -lt 2 ]; then
    echo "Usage: ./query.sh <channel-num> <function> [args...]"
    exit 1
fi

# Set env vars
source env-vars.sh

# Get arguments
function createArgs() {
    ARGS="\"$1\""
    for i in ${@:2}; do
        ARGS="$ARGS,\"$i\""
    done
    ARGS="{\"Args\":[$ARGS]}"
}

createArgs ${@:2}

# Query chaincode
peer chaincode query -C channel$1 -n basic$1 -c "$ARGS"