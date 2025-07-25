#!/bin/bash

if [ $# -lt 2 ]; then
    echo "Usage: ./invoke.sh <channel-num> <function> [args...] [--json]"
    exit 1
fi

# Set fabric env vars
source env-vars.sh

PEER_ORG_PATH="${PWD}/organizations/peerOrganizations"

# Check if last argument is --json
USE_JSON=false
if [[ "${!#}" == "--json" ]]; then
    USE_JSON=true
    # Remove the last argument (--json)
    set -- "${@:1:$(($#-1))}"
fi

function createCommand() {
    local FUNC_NAME=$1
    shift

    local ARGS_JSON=""

    for arg in "$@"; do
        if $USE_JSON; then
            # Pass argument as raw JSON (assumed valid)
            ARGS_JSON+="$arg,"
        else
            # Escape and quote argument as JSON string
            local arg_str=$(jq -c -R <<< "$arg")
            ARGS_JSON+="$arg_str,"
        fi
    done

    # Remove trailing comma
    ARGS_JSON="${ARGS_JSON%,}"

    COMMAND="{\"function\":\"$FUNC_NAME\",\"Args\":[$ARGS_JSON]}"
}

function createPeer0Connections() {
    for (( i=1; i<=$ORGANIZATION_NUMBER; i++ )); do
        local ORG_PATH="${PEER_ORG_PATH}/org${i}.example.com"
        local PEER_PATH="${ORG_PATH}/peers/peer0.org${i}.example.com"
        local PEER_TLS_CERT="${PEER_PATH}/tls/ca.crt"
        local PEER_ADDRESS="localhost:$((6 + i))050"
        PEER_CONNECTIONS="$PEER_CONNECTIONS --peerAddresses $PEER_ADDRESS --tlsRootCertFiles $PEER_TLS_CERT"
    done
}

createCommand "${@:2}"
createPeer0Connections

set -x
peer chaincode invoke \
    -o localhost:7000 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls \
    --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
    -C tradechannel$1 \
    -n traderchaincode$1 \
    $PEER_CONNECTIONS \
    -c "$COMMAND"
set +x
