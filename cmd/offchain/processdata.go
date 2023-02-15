package main

import "strconv"

func CreateAccount(AccountName string, SavingBalance int, CheckingBalance int) {
	// Create two separate accounts for two DB
	mSave[AccountName] = SavingBalance
	mCheck[AccountName] = CheckingBalance
}

func CreateAccountNum(accountNum int) {
	for i := 0; i < accountNum; i++ {
		CreateAccount(strconv.Itoa(i), 1000, 1000)
	}
}

func Dfs(GE [][]GraphEdge, groupNum int) ([]uint16, map[uint16]string) {
	sub := GE[groupNum]
	l := len(sub)
	m := make(map[uint16]string)
	s := NewSortedGraph(uint16(l))
	for i := 0; i < l; i++ {
		from := sub[i].F
		to := sub[i].T
		s.AddEdge(from, to)
		_, ok := m[to]
		if ok {
			m[to] += ">" + sub[i].D
		} else {
			m[to] = sub[i].D
		}
	}
	order := s.TopoSortByDFS()
	return order, m
}
