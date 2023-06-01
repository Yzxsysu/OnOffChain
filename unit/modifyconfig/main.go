package main

import (
	"fmt"
)

func main() {
	nodeNUm := []int{4, 7, 10, 13, 16, 19, 22, 25, 28}
	for _, n := range nodeNUm {
		for i := 0; i < n; i++ {
			filePath := fmt.Sprintf("/home/WorkPlace/github.com/Yzxsysu/onoffchain/testnodeconfig/%vnode/node%v/config/config.toml", n, i)
			fmt.Println(filePath)
			err := updateFileLine(filePath, i, n)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
