package main

import "fmt"

func main() {
	jobs := []int{2, 3, 4}
	results := make(chan int)

	for _, j := range jobs {
		go func(n int) {
			results <- n * n
		}(j)
	}

	sum := 0
	for i := 0; i < 4; i++ {
		sum += <-results
	}
	fmt.Println("二乗の合計:", sum)
}
