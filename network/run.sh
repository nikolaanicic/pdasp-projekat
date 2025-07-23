#!/bin/bash

./network.sh down
./network.sh up
./network.sh createChannel
./network.sh deployCC -ccp ../chaincode/ -ccn traderchaincode1 -c channel1
./network.sh deployCC -ccp ../chaincode/ -ccn traderchaincode2 -c channel2

# Init ledger
cd utils
./init_ledger.sh