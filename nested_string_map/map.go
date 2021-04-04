package nested_string_map

// T is an infinitely nested map with strings as keys
type T struct {
	items mapType
}

type mapType map[string]interface{}

// At queries the map_nested for values at the key
func (t T) At(path ...string) (v interface{}, found bool) {
	if t.items == nil {
		return nil, false
	}
	current := t.items
	var ok bool
	for i, k := range path {
		if v, ok = current[k]; !ok {
			return nil, false
		} else {
			if len(path)-1 == i {
				found = true
				break
			} else {
				current = v.(mapType)
			}
		}
	}
	// went through all the paths, but didn't get to an edge, so found "nothing"
	return
}

// Set assigns a value to a path and will create the intermediate nested maps in between for you.
func (t *T) Set(value interface{}, path ...string) {
	if t.items == nil {
		t.items = make(mapType)
	}
	var current mapType
	current = t.items
	var next mapType
	var ok bool
	for i, k := range path {
		if len(path)-1 == i {
			current[k] = value
		} else {
			if next, ok = current[k].(mapType); !ok {
				m := make(mapType)
				current[k] = m
				current = m
			} else {
				current = next
			}
		}
	}
}

func (t *T) Remove(path ...string) {
	if len(path) == 0 {
		// no path given, do nothing
		return
	}
	if t.items == nil {
		return
	}
	parentPath := path[0 : len(path)-1]
	if parent, ok := t.At(parentPath...); !ok {
		// not found, nothing to remove
		return
	} else {
		delete(parent.(mapType), path[len(path)-1])
	}
}
