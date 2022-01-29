package shoppinglist

import "testing"

func TestAddItem(t *testing.T) {
	shoppinglist := []string{}

	assertItems(t, AddItem(shoppinglist, "milk"), 1)
}

func assertItems(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("obtengo %d pero esperaba %d", got, want)
	}
}
