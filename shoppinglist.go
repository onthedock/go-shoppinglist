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

func RemoveItem(sl ShoppingList, item Item) int {
	for i, li := range sl {
		if li == item {
			sl[i] = sl[len(sl)-1]
			sl = sl[:len(sl)-1]
			return len(sl)
		}
	}
	return len(sl)
}
