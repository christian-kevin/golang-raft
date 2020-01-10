NODE_ID=$1

KV_PORT_NUMBER=1300${NODE_ID}

n=2

PEERS="http://127.0.0.1:12001"

while [[ $n -le $NODE_ID ]]
do
  PEERS="${PEERS},http://127.0.0.1:1200${n}"
  n=$((n+1))
done

COMMAND="./golang-raft --id ${NODE_ID} --cluster ${PEERS} --port ${KV_PORT_NUMBER} --join"

eval $COMMAND