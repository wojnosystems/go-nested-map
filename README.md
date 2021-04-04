# Overview

Creates a nesting set of maps with string keys that will store any type in the leaves of the map. The idea was to create a nil-coalesced structure in Go to avoid returning bad values if they do not exist along the path.

Paths can be any arbitrary slice of strings. The base type, nested_string_map.T will handle cleaning up memory for you to avoid memory leaks when you remove items from the list.

# Example

`go get github.com/wojnosystems/go-nested-map`

The below example creates a custom struct that only stores int values with string keys. This is a convenience method to avoid casting to int all over your code and I strongly recommend this practice for any type you wish to store in this structure.

```go
package main

import (
	"fmt"
	"github.com/wojnosystems/go-nested-map/nested_string_map"
	"strings"
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
	population := IntBag{}
	population.Set( 50_000, "usa", "ca", "yolo" )
	population.Set( 1_000_000, "usa", "ca", "alameda" )

	// what's the population of alameda county again?
	alamedaPopulation, ok := population.AtInt("usa", "ca", "alameda")
	if ok {
		fmt.Println("alameda:", alamedaPopulation)
	} else {
		fmt.Println("Hmmm... I don't know the population of alameda")
	}

	// what's the population of maricopa county again?
	maricopaPopulation, ok := population.AtInt("usa", "az", "maricopa")
	if ok {
		fmt.Println("maricopa:", maricopaPopulation)
	} else {
		fmt.Println("Hmmm... I don't know the population of maricopa")
	}

	fmt.Println("there are", population.Len(), "records")

	fmt.Println("I know about these counties:")
	for _, keys := range population.Keys() {
		fmt.Println(strings.Join(keys, ", "))
	}
}
```

outputs:

```text
alameda: 1000000
Hmmm... I don't know the population of maricopa
there are 2 records
I know about these counties:
usa, ca, yolo
usa, ca, alameda
```

# FAQ

## I want it to return a default value, though

Yea, these are really nice interfaces. However, since this underlying nested_string_map doesn't know what type, and we want to avoid casting whenever your custom type is used, this is left to you to implement. Instead of a method like: GetInt, you can write a method like "GetIntWithDefault(defaultValue int, path ...string) (value int)" and ensure that any ok = false cases are simply masked with the defaultValue that you wish to use.
