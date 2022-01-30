package main

import (
	"demoapp/shoppinglist"
	"fmt"
)

func PrintShoppingList(sl shoppinglist.ShoppingList) string {
	var l = shoppinglist.Item("Mi lista de la compra es:")
	for _, item := range sl {
		l += " " + item
	}
	return string(l)
}

func main() {
	sl := shoppinglist.ShoppingList{"milk", "sugar", "bread"}
	fmt.Println(PrintShoppingList(sl))
}
