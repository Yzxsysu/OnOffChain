GOSRC=/home/WorkPlace
ROOT=$GOSRC/github.com/Yzxsysu/OnOffChain

mkdir -p build

go build -o build/chain $ROOT/cmd/chain
go build -o build/offchain $ROOT/cmd/offchain
go build -o build/client $ROOT/cmd/client

chmod +x build/*
