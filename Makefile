.PHONY: build runOnChain runChain


# 定义传递的参数
ifndef groupNum
 groupNum = 4
endif
ifndef nodeId
  nodeId = 1
endif
ifndef division
  division = f
endif

# 输出参数
runChain:
	@echo "groupNum: $(groupNum)"
	@echo "nodeId: $(nodeId)"
	@echo "division: $(division)"
	bash ./scripts/$(groupNum)node/run$(groupNum)node.sh $(groupNum) $(nodeId) $(division)

build:
	bash ./scripts/goBuildExecute.sh

runOnChain:
	bash ./scripts/runOnChain.sh

