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

### Refactorizar - `ItemPresent()`

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

## Eliminando elementos de la lista de la compra

Como antes, primero vamos a introducir una función que nos permita eliminar un elemento (presente) en la lista de la compra.

### Creamos el test que valide la eliminación del elemento

Empezamos por definir el test:

```go
func TestRemoveItem(t *testing.T) {
    sl := ShoppingList{"milk", "sugar"}
    assertItems(t, RemoveItem(sl, "sugar"), 1)
}
```

### Eliminamos los errores de compilación

La ejecución del test vuelve a mostrar errores de compilación, ya que la función `RemoveItem` no existe todavía.

```go
func RemoveItem(sl ShoppingList, item Item) int {
    return 0
}
```

Usamos **el mínimo código posible** para eliminar los errores de compilación:

```bash
$ go test
--- FAIL: TestRemoveItem (0.00s)
    shoppinglist_test.go:26: obtengo 0 elementos en la lista pero esperaba 1
FAIL
exit status 1
FAIL    shoppinglist    0.003s
```

### Hacemos que la ejecución de los tests devuelva `PASS`

Ahora introducimos la mínima código posible para pasar el test:

```go
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
```

He usado el método *rápido*, que no preserva el orden de los elementos en el *slice* descrito en [2 ways to delete an element from a slice](https://yourbasic.org/golang/delete-element-slice/).

Este método sustituye el último elemento en el *slice* en la posición en la que hay una coincidencia. Esto elimina el elemento que queremos, pero deja un duplicado en las posiciones `i` y `len(slice)-1` (la última posición en el *slice*). Para eliminar el duplicado, copiamos todos los elementos **menos el último** a un nuevo *slice*, con lo que eliminamos el elemento duplicado.

Si el elemento `item` no está en la lista de la compra, no hacemos nada.

Validamos que el test pasa:

```bash
$ go test 
PASS
ok      shoppinglist    0.002s
```

Antes de pasar a refactorizar, quizás deberíamos añadir otro test para validar que en caso de que el elemento que se quiere eliminar no se encuentra en la lista todo funciona correctamente...

```go
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
```

### ¿Refactorizamos `RemoveItem`?

Al crear la función `ItemPresent` parecía que podríamos reusarla tanto al añadir como al eliminar un elemento de la lista... Pero acabamos de ver que para eliminar un elemento de un *slice* es necesario conocer la posición en la que se encuentra. La función `ItemPresent` sólo devuelve `true` si encuentra una coincidencia, pero no devuelve la posición.

Aquí es donde deberíamos aprovechar que en Go, las funciones pueden devolver múltiples valores. Lo *idiomático* en Go sería (supongo) devolver un entero (con la posición de la coincidencia, o -1, si no se encuentra) y un error (*nil* o un error, si no se ha encontrado el *item*).

Antes de empezar a modificar `RemoveItem`, adaptamos la función `ItemPresent`:

```go
func ItemPresent(sl ShoppingList, item Item) (int, error) {
    for i, li := range sl {
        if li == item {
            return i, nil
        }
    }
    return -1, errors.New("item not found")
}
```

Cambiamos el valor `bool` por `(int, error)`, para la posición en la que se ha encontrado la coincidencia y un error, en caso de no se haya encontrado.

Como antes, recorremos la lista de la compra, pero esta vez sí que estamos interesados en el valor del índice en el *slice*. Si el elemento se encuentra en la lista de la compra, devolvemos el índice y `nil`.

En el caso de que no se encuentre, devolvemos `-1` y un nuevo error indicando que no se ha encontrado el item en la lista de la compra.

A continuación tenemos que adaptar la función `AddItem`, donde se usa `ItemPresent`.

Llamamos a la función e ignoramos el índice (no nos interesa para la función `AddItem`). Si `ItemPresent` devuelve error, significa que no se ha encontrado el item a la lista de compra, y por tanto, lo añadimos.

Para finalizar, devolvemos la longitud de la lista de la compra.

Tras estas modificaciones, volvemos a verificar que los tests pasan:

```bash
$ go test
PASS
ok      shoppinglist    0.002s
```

El siguiente paso es modificar `RemoveItem` para usar también `ItemPresent`.

```go
func RemoveItem(sl ShoppingList, item Item) int {
    i, err := ItemPresent(sl, item)
    if err != nil {
        return len(sl)
    }

    sl[i] = sl[len(sl)-1]
    sl = sl[:len(sl)-1]
    return len(sl)
}
```

Comprobamos si el elemento `item` se encuentra en `sl`. Si no se encuentra (`err != nil`), no es necesario eliminarlo, así que devolvemos la longitud actual del *slice*.

En caso de que no haya error (y por tanto sí que exista el elemento en la lista de la compra), lo eliminamos y devolvemos la longitud del *slice* actualizado.

## Convertir las funciones en métodos

Podemos asociar una función a un tipo y convertirlo en un método. Dado que las funciones `AddItem` y `RemoveItem` actúan sobre una variable de tipo `ShoppingList`, lo ideal sería que formaran parte del propio tipo. Así podríamos llamarlas como `sl.Add` y `sl.Remove`, donde `sl` es una variable de tipo `ShoppingList`.

Empezamos modificando el test para la función de añadir un elemento a la lista de la compra:

```go
func TestAdd(t *testing.T) {
    t.Run("Add item to list", func(t *testing.T) {
        sl := ShoppingList{}
        assertItems(t, sl.Add("milk"), 1)
    })

    t.Run("Avoid adding duplicate item", func(t *testing.T) {
        sl := ShoppingList{"sugar"}
        assertItems(t, sl.Add("sugar"), 1)
    })
}
```

Al ejecutar `go test`:

```bash
$ go test
# shoppinglist [shoppinglist.test]
./shoppinglist_test.go:8:20: sl.Add undefined (type ShoppingList has no field or method Add)
FAIL    shoppinglist [build failed]
```

Solucionamos los problemas de compilación mediante:

```go
func (sl ShoppingList) Add(item Item) int {
    _, err := ItemPresent(sl, item)
    if err != nil {
        sl = append(sl, item)
    }
    return len(sl)
}
```

## Convertir las funciones en *métodos* para el tipo `ShoppingLIst`

Poco a poc estamos viendo que la funcionalidad a la lista de la compra la podemos encapsular en el nuevo tipo `ShoppingList`. Así que el siguiente paso natural es el de convertir las funciones en *métodos*.

Para ello, modificamos los tests asociados a la función `AddItem` que convertiremos en `Add`:

```go
func TestAdd(t *testing.T) {
    t.Run("Add item to list", func(t *testing.T) {
        sl := ShoppingList{}
        assertItems(t, sl.Add("milk"), 1)
    })

    t.Run("Avoid adding duplicate item", func(t *testing.T) {
        sl := ShoppingList{"sugar"}
        assertItems(t, sl.Add("sugar"), 1)
    })
}
```

A parte de moficicar el nombre del test, también cambiamos cómo llamamos a la función en `AssertItems(t, sl.Add("milk"), 1)`.

Para hacer solucionar los errores de compilación, cambiamos el nombre de la función de `AddItem` a `Add` y la *signature* de la función:

```go
func (sl ShoppingList) Add(item Item) int {
    _, err := ItemPresent(sl, item)
    if err != nil {
        sl = append(sl, item)
    }
    return len(sl)
}
```

Tras esta modificación, los tests pasan correctamente.

Así que realizamos la misma modificación para `RemoveItem` (que renombramos a `Remove`) y repetimos el proceso; modificamos el test:

```go
func TestRemoveItem(t *testing.T) {
    t.Run("Remove item", func(t *testing.T) {
        sl := ShoppingList{"milk", "sugar"}
        assertItems(t, sl.Remove("sugar"), 1)
    })
    t.Run("Do nothing if item is not found", func(t *testing.T) {
        sl := ShoppingList{"milk", "sugar"}
        assertItems(t, sl.Remove("bread"), 2)
    })
}
```

Esto provoca errores de compilación.

```bash
$ go test
# shoppinglist [shoppinglist.test]
./shoppinglist_test.go:27:20: sl.Remove undefined (type ShoppingList has no field or method Remove)
./shoppinglist_test.go:31:20: sl.Remove undefined (type ShoppingList has no field or method Remove)
FAIL    shoppinglist [build failed]
```

Los corregimos:

```go
func (sl ShoppingList) Remove(item Item) int {
    i, err := ItemPresent(sl, item)
    if err != nil {
        return len(sl)
    }

    sl[i] = sl[len(sl)-1]
    sl = sl[:len(sl)-1]
    return len(sl)
}
```

Y validamos que todo vuelve a estar OK:

```bash
$ go test
PASS
ok      shoppinglist    0.002s
```

## Convertir *ShoppingList* en una aplicación

Hemos definido un tipo `ShoppingList` y unos métodos para añadir o eliminar elementos a la lista de la compra.

El siguiente paso es usar estas funciones como parte de una aplicación que *haga algo*.

Para ello, movemos el fichero `shoppinglist.go` y `shoppinglist_test.go` a una carpeta llamada `shoppinglist`.

En la carpeta raíz, hemos inicializado el módulo llamado `demoapp` con `go mod init demoapp`.

Creamos un fichero para los tests `main_test.go` para seguir con la dinámica que hemos establecido.
La estructura de carpetas y ficheros queda:

```bash
$ tree 
.
├── go.mod
├── main_test.go
├── readme.md
└── shoppinglist
    ├── shoppinglist.go
    └── shoppinglist_test.go
```

### Crear el test

En el fichero `main_test.go` definimos el test para la aplicación que usará las funciones que hemos definido para gestionar la lista de la compra:

```go
package main

import (
    "demoapp/shoppinglist"
    "testing"
)

func TestPrintShoppingList(t *testing.T) {
    sl := shoppinglist.ShoppingList{"milk", "sugar"}
    got := PrintShoppingList(sl)
    want := "Mi lista de la compra es: milk sugar"
    if got != want {
        t.Errorf("obtengo %q pero quería %q", got, want)
    }
}
```

`shoppinglist` es un *package* del módulo `demoapp`; para poder usar las funciones definidas en el *package*, debemos importarlo. El nombre del *package* también es la ruta al *package*; en general, el nombre del módulo es la ruta al repositorio desde donde se puede obtener mediante `go get`, por lo que lo habitual es que sea de la forma `github.com/onthedock/demoapp`, por ejemplo.

Como para el resto de funciones, variables, etc, *importadas*, debemos precederlas del nombre del *package*.

Al ejecutar el test, obtenemos errores de compilación porque `PrintShoppingList` todavía no existe.

Lo creamos con el mínimo código posible para eliminar los errores de compilación:

```go
package main

import (
    "demoapp/shoppinglist"
    "fmt"
)

func PrintShoppingList(sl shoppinglist.ShoppingList) string {
    return ""
}

func main() {
    sl := shoppinglist.ShoppingList{"milk", "sugar", "bread"}
    fmt.Println(PrintShoppingList(sl))
}
```

Verificamos que los errores de compilación ya han sido solucionados:

```bash
$ go test
--- FAIL: TestPrintShoppingList (0.00s)
    main_test.go:13: obtengo "" pero quería "Mi lista de la compra es: milk sugar"
FAIL
exit status 1
FAIL    demoapp 0.006s
```

Ahora nos centramos en hacer que el test pase:

```go
func PrintShoppingList(sl shoppinglist.ShoppingList) string {
    var l = shoppinglist.Item("Mi lista de la compra es:")
    for _, item := range sl {
        l += " " + item
    }
    return string(l)
}
```

Y efectivamente, el test pasa:

```bash
$ go test
PASS
ok      demoapp 0.002s
```

## Conclusión

El procedimiento de crear primero el test, solucionar problemas de compilación, hacer que el test se verifique y refactorizar, una y otra vez, escribiendo en cada paso sólo la mínima cantidad de código permite avanzar de forma segura en el desarrollo de la aplicación. Cuando realizamos modificaciones obtenemos *feedback* inmediato -en forma de tests fallidos- si afectamos a funcionalidad existente (cuyos tests previamente habíamos validado).

Los tests nos ayudan a *pensar* en los detalles a implementar y nos ayudan a realizar modificaciones con confianza.

En esta aplicación *demo*, hemos seguido un camino inverso al que usaríamos habitualmente al construir una aplicación, en el que empezaríamos por un *package main* y cuando la aplicación alcanzara un tamaño poco manegable, la dividiríamos en módulos y *packages*. Como el objetivo era demostrar cómo usar el método de desarrollo basado en tests (y el resultado final de la "app" es el mismo), el orden en el que hemos desarrollado los *packages* no es relevante.
