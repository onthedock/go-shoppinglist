# ShoppingList: un ejercicio para aprender Go

Inspirado por [Learn Go with test](https://quii.gitbook.io/learn-go-with-tests/), voy a ir documentando el proceso de creación de una aplicación para gestionar una *lista de la compra* en Go.

## El primer test

Empezamos definiendo el primer test para validar que hemos añadido un elemento a la *lista de la compra*:

```go
package shoppinglist

import "testing"

func TestAddItem(t *testing.T) {
    shoppinglist := []string{}

    assertItems(t, AddItem(shoppinglist, "milk"), 1)
}

func assertItems(t *testing.T, got int, want int) {
    t.Helper()
    if got != want {
        t.Errorf("esperaba %d pero obtengo %d", got, want)
    }
}
```

Nuestra lista de la compra será un *slice* de `string`.

Verificaremos que se ha añadido un elemento a la lista de la compra comprobando la longitud del *slice*.

También definimos la función `func assertItems(t *testing.T, got int, want int) {}` como un *helper* (mediante `t.Helper()`) que nos permite evitar repetir la comprobación de `got` y `want` en cada test.

## Eliminando los problemas de compilación

En este estado, tenemos un problema a la hora de compilar la función:

```bash
$ go test
# shoppinglist [shoppinglist.test]
./shoppinglist_test.go:8:17: undefined: AddItem
FAIL    shoppinglist [build failed]
```

Como vemos, `AddItem` no está definido.

```go
package shoppinglist

func AddItem(shoppinglist []string, item string) int {
    return 0
}
```

Definimos `AddItem` de manera que acepta un *slice* de `string` y una `string` y devuelve un `int`. El mínimo código que necesitamos para la función es que ésta devuelva un `int`.

Ejecutando el test de nuevo comprobamos que ya no tenemos errores de compilación:

```bash
$ go test
--- FAIL: TestAddItem (0.00s)
    shoppinglist_test.go:8: obtengo 0 pero esperaba 1
FAIL
exit status 1
FAIL    shoppinglist    0.004s
```

## Refactorizar

Modificamos el código de la función `AddItem` para que el test pase.

El objetivo es escribir la menor cantidad de código para que el test pase; por ello, en este caso, lo único que haremos es añadir un elemento al *slice*.

```go
package shoppinglist

func AddItem(shoppinglist []string, item string) int {
    shoppinglist = append(shoppinglist, item)
    return len(shoppinglist)
}
```

Validamos ejecutando `go test`:

```bash
$ go test
PASS
ok      shoppinglist    0.002s
```

