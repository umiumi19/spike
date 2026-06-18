package main

import "fmt"

// sumNumbers は数値が入っているはずのスライスを合計する。
func sumNumbers(values []interface{}) int {
	total := 0
	for _, v := range values {
		total += v.(int)
	}
	return total
}

func main() {
	values := []interface{}{1, 2, "3", 4}
	fmt.Println("合計:", sumNumbers(values))
}
