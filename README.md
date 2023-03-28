# Introduction 
`claridate` provides several functions regarding the parsing of date and/or time.

# Getting Started
To start using this package, just import it and use the Functions like so:
```go
package main

import (
    "fmt"

    "github.com/Christoph-Harms/claridate"
)

func main() {
    dateFormat, err := claridate.DetermineDateFormat("2006-08")
    if err != nil {
        panic(err)
    }

    fmt.Println(dateFormat) // prints "YYYY-MM"
}
```
# Test
To run the tests, run `make test` in the project directory.

