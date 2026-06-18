package main

import "fmt"

// finalBalance は開始残高 start に、取引 deltas を順に適用した最終残高を返すつもり。
func finalBalance(start int, deltas []int) int {
	balance := start
	for _, d := range deltas {
		balance = start + d
	}
	return balance
}

func main() {
	start := 1000
	deltas := []int{-200, 500, -100, -50, 300}
	fmt.Println("最終残高:", finalBalance(start, deltas))
}
