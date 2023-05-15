package main

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/rand"
	"sync"
)

var sMap sync.Map

func main() {
	c := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		sMap = sync.Map{}
		go do(&sMap, c)
		go do1(&sMap, c)
		<-c
		<-c
		fmt.Printf("%p\n", &sMap)
	}
}

func do1(p *sync.Map, c chan bool) {
	r := rand.NewRand()
	rint := r.Intn(10)
	p.LoadOrStore("B", rint)
	c <- true
}

func do(p *sync.Map, c chan bool) {
	r := rand.NewRand()
	rint := r.Intn(10)
	p.LoadOrStore("A", rint)
	c <- true
}
