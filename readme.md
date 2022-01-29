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

### Eliminando los problemas de compilación

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

### Refactorizar

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

## Crear un tipo específico

Vamos a definir tipos específicos que esperamos que haga más sencillo de entender el código.

```go
type Item string
type ShoppingList []Item
```

Empezamos actualizando el test para definir `shoppinglist` de tipo `ShoppingList`:

```go
func TestAddItem(t *testing.T) {
    shoppinglist := ShoppingList{}

    assertItems(t, AddItem(shoppinglist, "milk"), 1)
}
```

Al ejecutar `go test`, encontramos errores de compilación:

```bash
$ go test
# shoppinglist [shoppinglist.test]
./shoppinglist_test.go:6:18: undefined: ShoppingList
FAIL    shoppinglist [build failed]
```

Vamos a definir los nuevos tipos (en `shoppinglist.go`):

```go
type Item string
type ShoppingList []Item
```

Tenemos que modificar la función `AddItem` para reflejar los nuevos tipos de los parámetros para la función:

```go
func AddItem(shoppinglist ShoppingList, item Item) int {
...
```

Una vez actualizado, validamos que el test sigue pasando:

```bash
$ go test
PASS
ok      shoppinglist    0.002s
```

## Sólo debemos añadir un nuevo elemento a la lista si no está ya en ella

La función `AddItem` añade un elemento a la lista de la compra tanto si el *item* ya está en ella como si no.

Vamos a añadir el requerimiento de que el elemento sólo debe añadirse si no está ya en la lista (no tiene sentido apuntar dos veces que tenemos que comprar *leche*, por ejemplo).

### Diseñamos un nuevo test

Antes de añadir un nuevo test para validar que no se añaden elementos que ya están presentes en la lista, convertimos el test existente en un [*subtest*](https://pkg.go.dev/testing#hdr-Subtests_and_Sub_benchmarks):

```go
func TestAddItem(t *testing.T) {
    t.Run("Add item to list", func(t *testing.T) {
        shoppinglist := ShoppingList{}
        assertItems(t, AddItem(shoppinglist, "milk"), 1)
    })
}
```

De esta forma podemos aplicar varios tests a la misma función.

El nuevo test queda:

```go
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
```

Ejecutando el test, vemos que falla:

```bash
$ go test
--- FAIL: TestAddItem (0.00s)
    --- FAIL: TestAddItem/Avoid_adding_duplicate_item (0.00s)
        shoppinglist_test.go:13: obtengo 2 pero esperaba 1
FAIL
exit status 1
FAIL    shoppinglist    0.003s
```

El mensaje del error podría mejorarse para indicar qué es lo que obtenemos y qué es lo que esperamos.

Actualizamos la funcion `assertItems`:

```bash
$ go test
--- FAIL: TestAddItem (0.00s)
    --- FAIL: TestAddItem/Avoid_adding_duplicate_item (0.00s)
        shoppinglist_test.go:13: obtengo 2 elementos en la lista pero esperaba 1
FAIL
exit status 1
FAIL    shoppinglist    0.003s
```

### Hacer que el test pase

Antes de añadir un elemento en la lista de la compra, tenemos que revisar si ya está en la lista.

Lo conseguimos recorriendo la *shoppinglist* y revisando si alguno de los elementos de la lista coincide con el nuevo elemento que queremos añadir:

```go
func AddItem(shoppinglist ShoppingList, item Item) int {
    for _, li := range shoppinglist {
        if li == item {
            return len(shoppinglist)
        }
    }
    shoppinglist = append(shoppinglist, item)
    return len(shoppinglist)
}
```

Validamos que los test pasan:

```bash
$ go test
PASS
ok      shoppinglist    0.002s
```

### Refactorizar - `IsItemPresent()`

Lo de tener que buscar si un elemento ya está en la lista de la compra será algo que tendremos que reutilizar (por ejemplo, cuando querramos eliminar un elemento de la lista).

De momento, lo convertimos en una función específica. Siguiendo con esa idea de usar el código más sencillo posible, esta nueva función devolverá `true` si ha encontrado el elemento y `false` en caso contrario.

Aprovechamos para *reducir* el nombre de la instancia de `ShoppingList` a `sl`. Dejamos `item` en vez de acortarlo a `i` para evitar confusiones con un *índice* de iteración en un bucle o similar.

```go
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
```

Validamos que tras la modificación los tests siguen pasando.

Si queremos ver el detalle de los tests (y los subtests), usamos `go test -v`:

```bash
$ go test -v
=== RUN   TestAddItem
=== RUN   TestAddItem/Add_item_to_list
=== RUN   TestAddItem/Avoid_adding_duplicate_item
--- PASS: TestAddItem (0.00s)
    --- PASS: TestAddItem/Add_item_to_list (0.00s)
    --- PASS: TestAddItem/Avoid_adding_duplicate_item (0.00s)
PASS
ok      shoppinglist    0.002s
```
