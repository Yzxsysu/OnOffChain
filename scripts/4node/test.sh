groupNum=$1
nodeId=$2
echo "group number: $groupNum, node id: $nodeId"
echo "cp -r /home/WorkPlace/github.com/Yzxsysu/onoffchain/testnodeconfig/${groupNum}node/* $TM_HOME"

case $nodeId in
    0)
    ./build/chain/chain -home $TM_HOME/${groupNum}node/node${nodeId} -leader "true" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 0 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6" -webPort "10157,10157,10157" -SetNum "2f" -subscribeIp "" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    1)
    ./build/chain/chain -home $TM_HOME/${groupNum}node/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 1 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6" -webPort "10157,10157,10157" -SetNum "2f" -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    2)
    ./build/chain/chain -home $TM_HOME/${groupNum}node/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 2 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6" -webPort "10157,10157,10157" -SetNum "2f" -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    3)
    ./build/chain/chain -home $TM_HOME/${groupNum}node/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 3 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6" -webPort "10157,10157,10157" -SetNum "2f"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
    -1)
    ./build/offchain/offchain -accountNum 1000 -webIp "127.0.0.1,127.0.0.1,127.0.0.1" -webPort "10157,10257,10357" -offChainIp "127.0.0.1" -offChainPort "8090" -SetNum "2f" &> $LOG_DIR/offchain.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
esac

case $nodeId in
    0)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "true" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 0 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "2f" -subscribeIp "" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    1)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 1 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "2f" -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    2)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 1 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "2f" -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    3)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 2 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "2f"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
    4)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 2 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "2f"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
    5)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 3 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157,10157,10157,10157" -SetNum "2f"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
    6)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -leaderIp "172.172.0.3:26657" -accountNum 1000 -OffChainIp "172.172.0.2" -OffChainPort "8090" -group 3 -coreNum 16 -webIp "172.172.0.4,172.172.0.5,172.172.0.6,172.172.0.7,172.172.0.8,172.172.0.9" -webPort "10157,10157,10157" -SetNum "2f"  -subscribeIp "127.0.0.1:10157" &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
    -1)
    ./build/offchain/offchain -accountNum 1000 -webIp "127.0.0.1,127.0.0.1,127.0.0.1" -webPort "10157,10257,10357" -offChainIp "127.0.0.1" -offChainPort "8090" -SetNum "2f" &> $LOG_DIR/offchain.log &
    echo "the node Id is ${nodeId}, offchain node"
    ;;
esac