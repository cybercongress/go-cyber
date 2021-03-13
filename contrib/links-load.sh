#!/bin/sh

CNT=0
ITER=$1
SLEEP=$2
NUMBLOCKS=$3
NODEADDR=$4
CHAINID=$5

if [ -z "$1" ]; then
  echo "Need to input number of iterations to run..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input number of seconds to sleep between iterations"
  exit 1
fi

if [ -z "$3" ]; then
  echo "Need to input block height to declare completion..."
  exit 1
fi

if [ -z "$4" ]; then
  echo "Need to input node address to poll..."
  exit 1
fi

if [ -z "$5" ]; then
  echo "Need to input chain-id..."
  exit 1
fi

#MNEMONIC="armed priority melt erode nurse profit drive mad major minimum border develop junk noodle empower orphan race girl tank result wealth cancel impulse hero"
#echo $MNEMONIC | cyber keys add tester --recover --keyring-backend test

while [ ${CNT} -lt $ITER ]; do
  curr_block=$(curl -s $NODEADDR:26657/status | jq -r '.result.sync_info.latest_block_height')

  if [ ! -z ${curr_block} ] ; then
    echo "Number of Blocks: ${curr_block}"
  fi

  if [ ! -z ${curr_block} ] && [ ${curr_block} -gt ${NUMBLOCKS} ]; then
    echo "Number of blocks reached. Success!"
    exit 0
  fi

#  if ! ((${CNT} % 5)); then
    R1=$RANDOM
    R2=$RANDOM
    FROM=$(echo $R1 | ipfs add -Q)
    TO=$(echo $R2 | ipfs add -Q)
    ADDR=$(cyber keys show tester3 -a --keyring-backend="test")
    cyber tx graph create $FROM $TO --from $ADDR --chain-id $CHAINID --keyring-backend="test" --node "tcp://$NODEADDR:26657" -b block -y
#  fi

  let CNT=CNT+1
  sleep $SLEEP
done

echo "Timeout reached. Failure!"
exit 1