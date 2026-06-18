package main

import "fmt"

type Sale struct {
	Product  string
	Revenue  int
	Quantity int
}

// averagePrice は1個あたりの平均単価を返す。
func averagePrice(s Sale) int {
	return s.Revenue / s.Quantity
}

func main() {
	sales := []Sale{
		{Product: "りんご", Revenue: 1000, Quantity: 10},
		{Product: "バナナ", Revenue: 2000, Quantity: 20},
		{Product: "さくらんぼ", Revenue: 500, Quantity: 0},
		{Product: "なつめ", Revenue: 3000, Quantity: 30},
	}

	for _, s := range sales {
		fmt.Printf("%s: 平均単価 = %d\n", s.Product, averagePrice(s))
	}
}
