package store

import (
	"fmt"
	"sync"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	mu    sync.Mutex
	items = []Item{{ID: "1", Name: "apple"}, {ID: "2", Name: "banana"}}
	seq   = 3
)

func List() []Item {
	mu.Lock()
	defer mu.Unlock()
	out := make([]Item, len(items))
	copy(out, items)
	return out
}

func Get(id string) (Item, bool) {
	mu.Lock()
	defer mu.Unlock()
	for _, it := range items {
		if it.ID == id {
			return it, true
		}
	}
	return Item{}, false
}

func Create(name string) Item {
	mu.Lock()
	defer mu.Unlock()
	it := Item{ID: fmt.Sprintf("%d", seq), Name: name}
	seq++
	items = append(items, it)
	return it
}
