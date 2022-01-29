package shoppinglist

type Item string
type ShoppingList []Item

func AddItem(sl ShoppingList, item Item) int {
	if ItemPresent(sl, item) {
		return len(sl)
	}
	sl = append(sl, item)
	return len(sl)
}

func ItemPresent(sl ShoppingList, item Item) bool {
	for _, li := range sl {
		if li == item {
			return true
		}
	}
	return false
}
