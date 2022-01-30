package main

import (
	"demoapp/shoppinglist"
	"testing"
)

func TestPrintShoppingList(t *testing.T) {
	sl := shoppinglist.ShoppingList{"milk", "sugar"}
	got := PrintShoppingList(sl)
	want := "Mi lista de la compra es: milk sugar"
	if got != want {
		t.Errorf("obtengo %q pero quería %q", got, want)
	}
}
