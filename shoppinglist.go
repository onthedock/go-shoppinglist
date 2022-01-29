package shoppinglist

type Item string
type ShoppingList []Item

func AddItem(shoppinglist ShoppingList, item Item) int {
	shoppinglist = append(shoppinglist, item)
	return len(shoppinglist)
}
