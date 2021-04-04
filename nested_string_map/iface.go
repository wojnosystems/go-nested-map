package nested_string_map

// Immutable methods do not alter the map's contents or keys
type Immutable interface {

	// Len is the number of items in the nested_string_map
	Len() int

	// IsEmpty returns true if there is nothing in the map, false otherwise
	IsEmpty() bool

	// Keys returns a list of keys to every valid item in the map. This could be expensive if the map has lots of keys and sub-keys
	Keys() [][]string
}

// ImmutableGetter is the only type-specific method that requires casting. This is not included in Immutable as custom sub-types are expected to make their own Immutable interface which includes the above Immutable interface, plus their own AtTYPE method for getting
type ImmutableGetter interface {
	// Fetches the leaf value at the path. Will return ok = false if not a leaf node (contains only a value and not another map) or if the path was not found
	Get(path ...string) (v interface{}, ok bool)
}

// Mutable methods alter the map
type Mutable interface {
	// Set assigns a value to a path
	Put(value interface{}, path ...string)

	// Remove deletes a value at the path
	Remove(path ...string)
}

type Mapper interface {
	ImmutableGetter
	Immutable
	Mutable
}
