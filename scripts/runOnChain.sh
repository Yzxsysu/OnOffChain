GOSRC=/home/WorkPlace
TEST_SCENE="OnChain"
TM_HOME="/home/.tendermint"
WORKSPACE="$GOSRC/github.com/Yzxsysu/onoffchain"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE"
DURATION=120

rm -rf $TM_HOME

mkdir -p $TM_HOME
mkdir -p $LOG_DIR

cp -r /home/WorkPlace/github.com/Yzxsysu/onoffchain/config/* $TM_HOME
echo "configs generated"

pkill -9 chain
pkill -9 offchain

./build/chain/chain -home $TM_HOME/4node/shard0/node0 -leader "true" -leaderIp "127.0.0.1:20057" -accountNum 1000 -OffChainIp "127.0.0.1" -OffChainPort "8090" -group 0 -coreNum 16 -webIp "127.0.0.1,127.0.0.1,127.0.0.1" -webPort "10157,10257,10357" -SetNum "2f" -subscribeIp "" &> $LOG_DIR/node0.log &

./build/chain/chain -home $TM_HOME/4node/shard0/node1 -leader "false" -leaderIp "127.0.0.1:20057" -accountNum 1000 -OffChainIp "127.0.0.1" -OffChainPort "8090" -group 1 -coreNum 16 -webIp "127.0.0.1,127.0.0.1,127.0.0.1" -webPort "10157,10257,10357" -SetNum "2f" -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node1.log &

./build/chain/chain -home $TM_HOME/4node/shard0/node2 -leader "false" -leaderIp "127.0.0.1:20057" -accountNum 1000 -OffChainIp "127.0.0.1" -OffChainPort "8090" -group 2 -coreNum 16 -webIp "127.0.0.1,127.0.0.1,127.0.0.1" -webPort "10157,10257,10357" -SetNum "2f" -subscribeIp "127.0.0.1:10257" &> $LOG_DIR/node2.log &

./build/chain/chain -home $TM_HOME/4node/shard0/node3 -leader "false" -leaderIp "127.0.0.1:20057" -accountNum 1000 -OffChainIp "127.0.0.1" -OffChainPort "8090" -group 3 -coreNum 16 -webIp "127.0.0.1,127.0.0.1,127.0.0.1" -webPort "10157,10257,10357" -SetNum "2f"  -subscribeIp "127.0.0.1:10357" &> $LOG_DIR/node3.log &

./build/offchain/offchain -accountNum 1000 -webIp "127.0.0.1,127.0.0.1,127.0.0.1" -webPort "10157,10257,10357" -offChainIp "127.0.0.1" -offChainPort "8090" -SetNum "2f" &> $LOG_DIR/offchain.log &


echo "testnet launched"
echo "running for ${DURATION}s..."
sleep $DURATION

pkill -9 chain
pkill -9 offchain
echo "all done"
