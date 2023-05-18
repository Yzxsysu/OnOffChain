GOSRC=/home/WorkPlace
TEST_SCENE="OnChain4Node"
TM_HOME="/home/.tendermint"
WORKSPACE="$GOSRC/github.com/Yzxsysu/onoffchain"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE"
DURATION=360

rm -rf $TM_HOME

mkdir -p $TM_HOME
mkdir -p $LOG_DIR

groupNum=$1
nodeId=$2
division=$3
echo "group number: $groupNum, node id: $nodeId, division : $division"

cp -r /home/WorkPlace/github.com/Yzxsysu/onoffchain/testnodeconfig/${groupNum}node/* $TM_HOME
echo "configs generated"

pkill -9 chain
pkill -9 offchain

case $nodeId in
    0)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "true" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 0 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "$division" -subscribeIp "" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    1)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 1 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "$division" -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    2)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 1 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "$division" -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    3)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 2 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "$division"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
    4)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 2 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "$division"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
    5)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 3 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "$division"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    6)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 3 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "$division"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    -1)
    ./build/offchain/offchain -accountNum 1000 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -offChainIp "127.0.0.1" -offChainPort "8090" -SetNum "$division" &> $LOG_DIR/offchain.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
esac

echo "testnet launched"
echo "running for ${DURATION}s..."
sleep $DURATION

pkill -9 chain
pkill -9 offchain
echo "all done"
