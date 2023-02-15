GOSRC=/home/WorkPlace
TEST_SCENE="OnChain"
TM_HOME="/home/.tendermint"
WORKSPACE="$GOSRC/github.com/Yzxsysu/OnOffChain"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE"
DURATION=60

rm -rf $TM_HOME

mkdir -p $TM_HOME
mkdir -p $LOG_DIR

cp -r /home/WorkPlace/github.com/Yzxsysu/OnOffChain/config/* $TM_HOME
echo "configs generated"

pkill -9 chain
pkill -9 offchain

./build/chain/chain -home $TM_HOME/node/node0 -leader "true" -leaderIp "0.0.0.0:26657" -accountNum 1000 -OffChainIp "0.0.0.0:8090" -group 1 -coreNum 16 -SetNum "2f"  &> $LOG_DIR/node0.log &

./build/chain/chain -home $TM_HOME/node/node1 -leader "true" -leaderIp "0.0.0.0:26657" -accountNum 1000 -OffChainIp "0.0.0.0:8090" -group 1 -coreNum 16 -SetNum "2f"  &> $LOG_DIR/node1.log &

./build/chain/chain -home $TM_HOME/node/node2 -leader "true" -leaderIp "0.0.0.0:26657" -accountNum 1000 -OffChainIp "0.0.0.0:8090" -group 2 -coreNum 16 -SetNum "2f"  &> $LOG_DIR/node2.log &

./build/chain/chain -home $TM_HOME/node/node3 -leader "true" -leaderIp "0.0.0.0:26657" -accountNum 1000 -OffChainIp "0.0.0.0:8090" -group 3 -coreNum 16 -SetNum "2f"  &> $LOG_DIR/node3.log &

./build/offchain/offchain -leaderIp "0.0.0.0:26657" -accountNum 1000  &> $LOG_DIR/offchain.log &


echo "testnet launched"
echo "running for ${DURATION}s..."
sleep $DURATION

pkill -9 chain
pkill -9 offchain
echo "all done"
