package application

import "crypto/sha256"

func AddComplexity(byteLen, cycleNum int) {
	for i := 0; i < cycleNum; i++ {
		data := make([]byte, byteLen)
		//data := []byte("hello world, hello world, hello world, hello world")
		h := sha256.New()
		h.Write(data)
	}
}
