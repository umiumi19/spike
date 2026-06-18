package main

import "fmt"

// lastValue はスライスの末尾の要素を返すつもりの関数。
func lastValue(nums []int) int {
	return nums[len(nums)]
}

func main() {
	data := []int{10, 20, 30}
	fmt.Println("最後の値:", lastValue(data))
}
