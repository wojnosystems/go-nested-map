package main

import (
	"fmt"
	"github.com/wojnosystems/go-nested-map/nested_string_map"
	"strings"
)

type IntBag struct {
	nested_string_map.T
}

// GetInt only returns values if they exist and are ints, handles casting for the underlying nested_map struct
func (b *IntBag) GetInt(path ...string) (value int, ok bool) {
	v, ok := b.Get(path...)
	if ok {
		value, ok = v.(int)
	}
	return
}

type Immutable interface {
	nested_string_map.Immutable
	GetInt(path ...string) (value int, ok bool)
}

func main() {
	population := IntBag{}
	population.Put(50_000, "usa", "ca", "yolo")
	population.Put(1_000_000, "usa", "ca", "alameda")

	// what's the population of alameda county again?
	alamedaPopulation, ok := population.GetInt("usa", "ca", "alameda")
	if ok {
		fmt.Println("alameda:", alamedaPopulation)
	} else {
		fmt.Println("Hmmm... I don't know the population of alameda")
	}

	// what's the population of maricopa county again?
	maricopaPopulation, ok := population.GetInt("usa", "az", "maricopa")
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
