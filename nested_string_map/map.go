package nested_string_map

// T is an infinitely nested map with strings as keys
type T struct {
	// items is the root of the tree of all maps
	items mapType

	// itemCount is the number of items in all sub-trees of items
	itemCount int
}

type mapType map[string]interface{}

// At queries the map_nested for values at the key
func (t T) Get(path ...string) (v interface{}, found bool) {
	v, found = t.getWithTreeNodes(path...)
	if _, isMapType := v.(mapType); isMapType {
		found = false
	}
	return
}

// getWithTreeNodes is the core of Get. Returns the item or branch node of the path. Returns found = true if and only if path maps to a value or a branch node in the tree. v is undefined when found is false.
func (t T) getWithTreeNodes(path ...string) (v interface{}, found bool) {
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

func (t *T) Len() int {
	return t.itemCount
}

func (t *T) IsEmpty() bool {
	return t.Len() == 0
}

// Set assigns a value to a path and will create the intermediate nested maps in between for you.
func (t *T) Put(value interface{}, path ...string) {
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
			t.itemCount++
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
	if parent, ok := t.getWithTreeNodes(parentPath...); !ok {
		// not found, nothing to remove
		return
	} else {
		delete(parent.(mapType), path[len(path)-1])
		t.itemCount--
		t.deleteEmptyParentPaths(parentPath...)
	}
}

// deleteEmptyParentPaths goes up the tree from a starting point, deleting any branch nodes that are empty
// this reclaims any map memory that is no longer used to prevent memory leaks
func (t *T) deleteEmptyParentPaths(path ...string) {
	if t.items == nil {
		// this should never occur, as deleteEmptyPaths is only called if an item was deleted,
		// which means t.items was initialized
		return
	}
	stack := make([]mapType, 0, len(path))
	stack = append(stack, t.items)
	for _, key := range path {
		if branch, ok := last(stack)[key]; !ok {
			// didn't find the branch... stop
			break
		} else {
			// found the branch, traverse
			stack = append(stack, branch.(mapType))
		}
	}

	// parent paths collected
	for i := len(stack) - 1; i > 0; i-- {
		branch := stack[i]
		if len(branch) == 0 {
			// this branch is empty, remove it
			delete(stack[i-1], path[i-1])
		}
	}
}

// last returns the last item in the stack, convenience method to keep things cleaner looking
func last(stack []mapType) mapType {
	return stack[len(stack)-1]
}

func (t *T) Keys() (paths [][]string) {
	if t.items == nil {
		// no keys
		return [][]string{}
	}
	paths = make([][]string, 0, t.Len())
	keysDFSRecursive(t.items, []string{}, &paths)
	return
}

// keysDFSRecursive collects all of the keys for all leaf nodes in the map
// this method is not very efficient and could be build to be faster.
func keysDFSRecursive(m mapType, parentPath []string, output *[][]string) {
	for key, child := range m {
		path := make([]string, len(parentPath)+1)
		for i, k := range parentPath {
			path[i] = k
		}
		path[len(path)-1] = key
		if childT, ok := child.(mapType); ok {
			keysDFSRecursive(childT, path, output)
		} else {
			*output = append(*output, path)
		}
	}
}
