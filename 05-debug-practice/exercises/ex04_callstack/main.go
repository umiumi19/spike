package main

import "fmt"

// indexOf は items の中での target の位置を返す。見つからなければ -1。
func indexOf(items []string, target string) int {
	for i, it := range items {
		if it == target {
			return i
		}
	}
	return -1
}

// priceOf は商品名から価格を引く。
func priceOf(catalog []int, items []string, name string) int {
	idx := indexOf(items, name)
	return catalog[idx]
}

func main() {
	items := []string{"ペン", "ノート", "ランプ"}
	catalog := []int{120, 980, 3400}

	order := []string{"ノート", "ペン", "マグカップ"}

	total := 0
	for _, name := range order {
		total += priceOf(catalog, items, name)
	}
	fmt.Println("合計:", total)
}
