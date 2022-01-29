package shoppinglist

import "testing"

func TestAddItem(t *testing.T) {
	t.Run("Add item to list", func(t *testing.T) {
		sl := ShoppingList{}
		assertItems(t, AddItem(sl, "milk"), 1)
	})

	t.Run("Avoid adding duplicate item", func(t *testing.T) {
		sl := ShoppingList{"sugar"}
		assertItems(t, AddItem(sl, "sugar"), 1)
	})
}

func assertItems(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("obtengo %d elementos en la lista pero esperaba %d", got, want)
	}
}

func TestRemoveItem(t *testing.T) {
	t.Run("Remove item", func(t *testing.T) {
		sl := ShoppingList{"milk", "sugar"}
		assertItems(t, RemoveItem(sl, "sugar"), 1)
	})
	t.Run("Do nothing if item is not found", func(t *testing.T) {
		sl := ShoppingList{"milk", "sugar"}
		assertItems(t, RemoveItem(sl, "bread"), 2)
	})
}
