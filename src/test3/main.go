package main

import "fmt"

func main() {
	u := make([][]uint16, 0)
	a := make([]uint16, 10)
	u = append(u, a)
	u = append(u, u[0])
	fmt.Println(u)
}
