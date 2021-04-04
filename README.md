# Overview

Creates a nesting set of maps with string-base keys that will store any type in the leaves of the map. The idea was to create a
nil-coalesced structure in Go to avoid returning bad values if they do not exist along the path.

# Example

`go get github.com/wojnosystems/go-nested-map`

```go
package main

import (
	"fmt"
	"github.com/wojnosystems/go-nested-map/nested_string_map"
)

type IntBag struct{
	nested_string_map.T
}

// AtInt only returns values if they exist and are ints, handles casting for the underying nested_map struct
func (b *IntBag) AtInt(path ...string) (value int, ok bool) {
	v, ok := b.At(path...)
	if ok {
		value, ok = v.(int)
	}
	return
}

func main() {
	population := IntBag{T: nested_string_map.New()}
	population.Set( 50_000, "usa", "ca", "yolo" )
	population.Set( 1_000_000, "usa", "ca", "alameda" )

	// what's the population of alameda county again?
	alamedaPopulation, ok := population.AtInt("usa", "ca", "alameda")
	if ok {
		fmt.Printf(`alameda: %d\n`, alamedaPopulation)
	} else {
			fmt.Println("Hmmm... I don't know the population of alameda")
		}

	// what's the population of maricopa county again?
	maricopaPopulation, ok := population.AtInt("usa", "az", "maricopa")
	if ok {
		fmt.Printf(`maricopa: %d\n`, maricopaPopulation)
	} else {
		fmt.Println("Hmmm... I don't know the population of maricopa")
    }
}
```

outputs:

```text
alameda: 1000000
Hmmm... I don't know the population of maricopa
```
