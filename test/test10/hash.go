package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	/*for i := 0; i < 1; i++ {
		data := make([]byte, 64)
		//data := []byte("hello world, hello world, hello world, hello world")
		h := sha256.New()
		h.Write(data)
		hash := h.Sum(nil)
		fmt.Printf("%x", hash)
	}*/
	AddComplexity(0, 0)
	elapsed := time.Since(start)
	fmt.Printf("代码执行时间为：%v\n", elapsed)
}

func AddComplexity(byteLen, cycleNum int) {
	for i := 0; i < cycleNum; i++ {
		data := make([]byte, byteLen)
		//data := []byte("hello world, hello world, hello world, hello world")
		h := sha256.New()
		h.Write(data)
	}
}

func int64ToBytes(n int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, n)
	return buf
}
