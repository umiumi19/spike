package main

import "testing"

func TestFinalBalance(t *testing.T) {
	got := finalBalance(1000, []int{-200, 500, -100, -50, 300})
	want := 1450
	if got != want {
		t.Errorf("finalBalance() = %d, want %d", got, want)
	}
}
