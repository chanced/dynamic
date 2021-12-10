# dynamic

dynamic is a collection of dynamic data types to supoprt [picker](https://github.com/chanced/picker)

All types have a `Value` method which return `interface{}`. The reason for this is that they can be `nil` or potentially various types. Values that are `nil` are json encoded to `null`. Use a pointer if you avoid this as json `omitempty` will not work.

dynamic types are not thread safe.

## dynamic.Bool

You can set Bool with any of the following:

-   `bool`, `*bool`
-   `dynamic.Bool`
-   `*dynamic.Bool`
-   `string`
-   `*string`
-   `[]byte`
-   `fmt.Stringer`
-   `nil`

```go
package main
import (
    "fmt"
    "github.com/chanced/dynamic"
)

boolean, err := dynamic.NewBool("true")
_ = err
if b, ok := boolean.Bool(); ok {
    fmt.Println(b)
}
err = boolean.Set("false")
_ = err

boolean, err = dynamic.NewBool(true)
_ = err // handle err
if b, ok := boolean.Bool(); ok {
    fmt.Println(b)
}
err = boolean.Set("1")
_ = err
```

## dynamic.Number

You can set Number with any of the following:

-   `string`,
-   `*string`
-   `json.Number`
-   `fmt.Stringer`
-   `int`,
-   `int`, `int64`, `int32`, `int16`, `int8`
-   `uint`, `uint64`, `uint32`, `uint16`, `uint8`
-   `float64`, `float32`
-   `complex128`, `complex64`
-   `*int`, `*int64`, `*int32`, `*int16`, `*int8`
-   `*uint`, `*uint64`, `*uint32`, `*uint16`, `*uint8`
-   `*float64`, `*float32`,
-   `*complex128`, `*complex64`
-   `[]byte`,
-   `fmt.Stringer`
-   `nil`

```go
package main
import (
    "fmt"
    "github.com/chanced/dynamic"
    "math"
)
func main(){
    number, err := dynamic.NewNumber(34)
    _ = err
    if u, ok := number.Uint64(); ok {
        fmt.Println(u)
    }
    err = n.Set("34.34")
    _ = err
    if u, ok := number.Uint64(); ok {
        // this wont be reached because 34.34 can not be
        // converted to a float without losing data
        fmt.Println(u)
    }
    if f, ok := number.Float32(); ok {
      fmt.Println(f)
    }

    err = number.Set(math.MaxFloat64)
    _ = err // no err but demonstrating

    if f, ok := number.Float32(); ok {
      // this won't be reached because number exceeds
      // MaxFloat32
    }
}
```

## dynamic.String

String accepts any of the following types:

-   `string`, `*string`
-   `[]byte`
-   `dynamic.String`, `*dynamic.String`
-   `fmt.Stringer`
-   `[]string` (joined with `","`)
-   `int`, `int64`, `int32`, `int16`, `int8`, `*int`, `*int64`, `*int32`, `*int16`, `*int8`,
-   `uint`, `uint64`, `uint32`, `uint16`, `uint8`, `*uint`, `*uint64`, `*uint32`, `*uint16`, `*uint8`
-   `float64`, `float32`, `complex128`, `complex64`, `*float64`, `*float32`, `*complex128`, `*complex64`
-   `bool`, `*bool`
-   `nil`

```go
package main
import (
    "fmt"
    "github.com/chanced/dynamic"
)

func main() {
    str, err := dynamic.NewString("str")
    val := str.ToLower().String() // "str"

    if str.Equal("true") {
        fmt.Println("equal")
    }
    fmt.Println(str.ToLower().String())
}

```

## dynamic.StringNumberBoolOrTime

StringNumberBoolOrTime accepts any of the following types:

-   `string`, `*string`
-   `[]byte`
-   `time.Time`, `*time.Time`
-   `fmt.Stringer`
-   `[]string` (joined with `","`)
-   `int`, `int64`, `int32`, `int16`, `int8`, `*int`, `*int64`, `*int32`, `*int16`, `*int8`,
-   `uint`, `uint64`, `uint32`, `uint16`, `uint8`, `*uint`, `*uint64`, `*uint32`, `*uint16`, `*uint8`
-   `float64`, `float32`, `complex128`, `complex64`, `*float64`, `*float32`, `*complex128`, `*complex64`
-   `bool`, `*bool`
-   `nil`

```go
package main
import (
    "fmt"
    "time"
    "github.com/chanced/dynamic"
)
func main() {
    now := dynamic.NewStringNumberBoolOrTime(time.Now())

    if n, ok := now.Time(); {
        fmt.Println(n)
    }
}
```

## dynamic.StringNumberOrTime

StringNumberOrTime accepts any of the following types:

-   `string`, `*string`
-   `[]byte`
-   `time.Time`, `*time.Time`
-   `fmt.Stringer`
-   `[]string` (joined with `","`)
-   `int`, `int64`, `int32`, `int16`, `int8`, `*int`, `*int64`, `*int32`, `*int16`, `*int8`,
-   `uint`, `uint64`, `uint32`, `uint16`, `uint8`, `*uint`, `*uint64`, `*uint32`, `*uint16`, `*uint8`
-   `float64`, `float32`, `complex128`, `complex64`, `*float64`, `*float32`, `*complex128`, `*complex64`
-   `bool`, `*bool`
-   `nil`

```go
package main
import (
    "log"
    "github.com/chanced/dynamic"
)
func main() {
    now := dynamic.NewStringNumberOrTime(34.34)

    if n, ok := now.Float64(); {
        log.Println(n)
    }

    if n, ok := now.Int64(); {
        // this won't be reached because it isn't possible to
        // convert the float value without losing data
        log.Fatal("was able to cast a float as int64")
    }
}
```

## dynamic.BoolOrString

BoolOrString accepts any of the following types:

-   `bool`, `*bool`
-   `string`, `*string`
-   `[]byte`
-   `fmt.Stringer`
-   `[]string` (joined with `","`)
-   `int`, `int64`, `int32`, `int16`, `int8`, `*int`, `*int64`, `*int32`, `*int16`, `*int8`,
-   `uint`, `uint64`, `uint32`, `uint16`, `uint8`, `*uint`, `*uint64`, `*uint32`, `*uint16`, `*uint8`
-   `float64`, `float32`, `complex128`, `complex64`, `*float64`, `*float32`, `*complex128`, `*complex64`
-   `nil`

```go
package main
import (
    "fmt"
    "time"
    "github.com/chanced/dynamic"
)
func main() {
    v := dynamic.NewBoolOrString("true")

    if n, ok := v.Bool(); {
        fmt.Println(n)
    }
}
```

## dynamic.StringOrArrayOfStrings

This is essentially a `[]string` except it'll Unmarshal either a `string` or a `[]string`. It always marshals into `[]string` though.

`NewStringOrArrayOfStrings` and `Set` accept:

-   `string`, `*string`
-   `[]byte`
-   `dynamic.String`, `*dynamic.String`
-   `fmt.Stringer`
-   `[]string` (joined with `","`)
-   `int`, `int64`, `int32`, `int16`, `int8`, `*int`, `*int64`, `*int32`, `*int16`, `*int8`,
-   `uint`, `uint64`, `uint32`, `uint16`, `uint8`, `*uint`, `*uint64`, `*uint32`, `*uint16`, `*uint8`
-   `float64`, `float32`, `complex128`, `complex64`, `*float64`, `*float32`, `*complex128`, `*complex64`
-   `bool`, `*bool`
-   `nil`

or you can use it like a slice:

```go
package main

import (
  "github.com/chanced/dynamic"
  "log"
)
func main() {
    strs := dynamic.StringOrArrayOfStrings{"value", "value2"}

    err := strs.Iterate(func(v string) error{
        if v == "value" {
            return dynamic.Done
        }
    })

    if err != nil {
        // err will be nil because Iterate checks for dynamic.Done
        // (or an error which returns "done" from Error())
        log.fatal(err)
    }

    err = strs.Iterate(func(v string) error) {
        if v == "value" {
            return fmt.Errorf("some error")
        }
    })
    if err != nil {
        // err will be "some error"
        log.Println(err)
    }

}

```

## dynamic.JSON

JSON is basically `[]byte` with helper methods as well as satisfying `json.Marshaler` and `json.Unmarshaler`

```go
  import(
      "fmt"
      "encoding/json"
      "github.com/chanced/dynamic"
  )

  func main() {
       data, _ := json.Marshal("str")
      d := dynamic.JSON(data)
      fmt.Println(d.IsString()) // true
      fmt.Println(d.IsBool()) // false


      // dynamic.JSON does not parse strings for potential
      // values:
      data, _ = json.Marshal("true")
      d = dynamic.JSON(data)
      fmt.Println(d.IsString()) // true
      fmt.Println(d.IsBool()) // false

      fmt.Println(d.UnquotedString()) // prints true

      data, _ = json.Marshal(true)
      d = dynamic.JSON(data)
      fmt.Println(d.IsString()) // false
      fmt.Println(d.IsBool()) // true
      fmt.Println(d.IsNumber()) // false
      fmt.Println(d.IsObject()) // false
      fmt.Println(d.IsArray()) // false

      data, = json.Marshal(map[string]string{"key":"value"})
      d = dynamic.JSON(data)
      fmt.Println(d.IsObject()) // true
      fmt.Println(d.IsEmptyObject()) // false

      data, = json.Marshal(map[string]string{"key":"value"})
      fmt.Println(d.IsObject()) // true
      fmt.Println(d.IsEmptyObject()) // true
  }
```

## Other types and mentions:

### dynamic.JSONObject

`dynamic.JSONObject` is a `map[string]dynamic.JSON`

### dynamic.Null

`dynamic.Null` is `[]byte("null")` for json purposes

### dynamic.Done

`dynamic.Done` is an `error` that indicates an iterator should stop but not return the error to the caller.

## TODO

-   [ ] Add `math` functions as methods to `Number`
-   [ ] Add `dynamic.String` methods to all types which could be `string`
-   [ ] Lot more testing to do
-   [ ] Comments

## License

Apache 2.0
