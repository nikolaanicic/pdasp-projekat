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
    json_mode=$2   
    shift 2
    args=("$@")   

    infoln "Testing invoke $func with mode $json_mode"


    if [[ "$json_mode" == "json" ]]; then

		local ARGS_JSON=""

		for arg in "$@"; do
			local arg_str=$(jq -c -R <<< "$arg")
			ARGS_JSON+="$arg_str,"
		done
		ARGS_JSON="${ARGS_JSON%,}"

		OUTPUT=$(./invoke.sh 1 "$func" "$ARGS_JSON" --json 2>&1)
    else
        OUTPUT=$(./invoke.sh 1 "$func" "${args[@]}" 2>&1)
    fi
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




infoln "Testing users"
invoke_function DeleteUser raw u1
USER_JSON='{"id":"u1","name":"Alice","last_name":"Alicee","email":"a@gmail.com","receipts_ids":[],"account_balance":100}'
invoke_function CreateUser json "$USER_JSON"
query_function ReadUser u1

query_function GetAllUsers


infoln "Testing rich queries"
QUERY_JSON='{"price":"2"}'
invoke_function QueryProducts json "$QUERY_JSON"  

query_function GetUsersGTEBalance 100
query_function SearchUsersByLastName Alicee
query_function SearchUsersByName Alice


infoln "Testing products"
invoke_function DeleteProduct raw pppp1
PRODUCT_JSON='{"id":"pppp1","name":"p1","expiration_date":"", "price":2,"quantity":2,"trader_id":"tt1"}'
invoke_function CreateProduct json "$PRODUCT_JSON"
query_function ReadProduct pppp1

invoke_function BuyProduct raw pppp1 u1
query_function GetAllProducts
query_function ReadUser u1



infoln "Testing traders"
invoke_function DeleteTrader raw tt111

TRADER_JSON='{"id":"tt111","trader_type":"MARKET","pib":"pppiiiibbb1","products":["br1"],"receipts":[],"account_balance":100}'
invoke_function CreateTrader json "$TRADER_JSON"

query_function GetAllTraders
query_function ReadTrader tt111

infoln "Testing receipts"
query_function GetAllReceips
 
successln "Success: $successNum"
warnln "Failed: $failNum"