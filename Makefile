.PHONY: build runOnChain

build:
	bash ./scripts/goBuildExecute.sh

runOnChain:
	bash ./scripts/runOnChain.sh

run_elrondtest:
	bash ./scripts/run_elrond.sh

run_haechitest:
	bash ./scripts/run_haechi.sh
