package shoppinglist

import (
	"errors"
)

type Item string
type ShoppingList []Item

func (sl ShoppingList) Add(item Item) int {
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

func (sl ShoppingList) Remove(item Item) int {
	i, err := ItemPresent(sl, item)
	if err != nil {
		return len(sl)
	}

	sl[i] = sl[len(sl)-1]
	sl = sl[:len(sl)-1]
	return len(sl)
}
