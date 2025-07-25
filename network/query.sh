#!/bin/bash

if [ $# -lt 2 ]; then
    echo "Usage: ./query.sh <channel-num> <function> [args...]"
    exit 1
fi

# Set env vars
source env-vars.sh



source env-vars.sh

PEER_ORG_PATH="${PWD}/organizations/peerOrganizations"

function createPeer0Connections() {
    for (( i=1; i<=$ORGANIZATION_NUMBER; i++ )); do
        local ORG_PATH="${PEER_ORG_PATH}/org${i}.example.com"
        local PEER_PATH="${ORG_PATH}/peers/peer0.org${i}.example.com"
        local PEER_TLS_CERT="${PEER_PATH}/tls/ca.crt"
        local PEER_ADDRESS="localhost:$((6 + $i))050"
        PEER_CONNECTIONS="$PEER_CONNECTIONS --peerAddresses $PEER_ADDRESS --tlsRootCertFiles $PEER_TLS_CERT"
    done
}

# Get arguments
function createArgs() {
    ARGS="\"$1\""
    for i in ${@:2}; do
        ARGS="$ARGS,\"$i\""
    done
    ARGS="{\"Args\":[$ARGS]}"
}



createArgs ${@:2}
createPeer0Connections

set -x
peer chaincode query -C tradechannel$1 -n traderchaincode$1 -c "$ARGS"

set +x