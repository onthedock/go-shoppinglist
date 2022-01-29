package shoppinglist

func AddItem(shoppinglist []string, item string) int {
	shoppinglist = append(shoppinglist, item)
	return len(shoppinglist)
}
