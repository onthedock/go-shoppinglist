package shoppinglist

import (
	"errors"
)

type Item string
type ShoppingList []Item

func AddItem(sl ShoppingList, item Item) int {
	_, err := ItemPresent(sl, item)
	if err != nil {
		sl = append(sl, item)
	}
	return len(sl)
}

func ItemPresent(sl ShoppingList, item Item) (int, error) {
	for i, li := range sl {
		if li == item {
			return i, nil
		}
	}
	return -1, errors.New("item not found")
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
