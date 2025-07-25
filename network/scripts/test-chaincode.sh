. scripts/utils.sh

successNum=0
failNum=0

function query_function() {
    func=$1
    args=${@:2}
    infoln "Testing query $func"
    ./query.sh 1 $func $args | json_pp
    if [ $? -ne 0 ]; then
        warnln "Failed to query $func"
        failNum=$(($failNum + 1))
    else
        successln "Test passed"
        successNum=$(($successNum + 1))
    fi
     
    echo ""
    sleep 3
}

function invoke_function() {
    func=$1
    args=${@:2}
    infoln "Testing invoke $func"

    OUTPUT=$(./invoke.sh 1 "$func" "$args" 2>&1)
    echo "$OUTPUT"

    if echo "$OUTPUT" | grep -q "Error: "; then
        warnln "Failed to invoke $func"
        failNum=$((failNum + 1))
    else
        successln "Test passed"
        successNum=$((successNum + 1))
    fi

    echo ""
    sleep 3
}

successln "Testing chaincode on tradechannel1"
sleep 1

infoln "Testing rich queries"


infoln "Testing users"
USER_JSON='{"id":"u1","name":"Alice"}'
invoke_function CreateUser "$USER_JSON"



successln "Success: $successNum"
warnln "Failed: $failNum"