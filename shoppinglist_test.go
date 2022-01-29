package shoppinglist

import "testing"

func TestAddItem(t *testing.T) {
	t.Run("Add item to list", func(t *testing.T) {
		shoppinglist := ShoppingList{}
		assertItems(t, AddItem(shoppinglist, "milk"), 1)
	})

	t.Run("Avoid adding duplicate item", func(t *testing.T) {
		shoppinglist := ShoppingList{"sugar"}
		assertItems(t, AddItem(shoppinglist, "sugar"), 1)
	})
}

func assertItems(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("obtengo %d elementos en la lista pero esperaba %d", got, want)
	}
}
