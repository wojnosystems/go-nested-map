package nested_string_map

var RootPath []string

// T is an infinitely nested map with strings as keys
type T struct {
	// items is the root of the tree of all maps
	items mapWithValue

	// itemCount is the number of items in all sub-trees of items
	itemCount int
}

type mapType map[string]mapWithValue

type mapWithValue struct {
	nested mapType
	value  interface{}
}

// At queries the map_nested for values at the key
func (t T) Get(path ...string) (v interface{}, found bool) {
	var branch mapWithValue
	branch, found = t.getWithTreeNodes(path...)
	if !found {
		return
	}
	if branch.value != nil {
		found = true
		v = branch.value
	} else {
		found = false
	}
	return
}

// getWithTreeNodes is the core of Get. Returns the item or branch node of the path. Returns found = true if and only if path maps to a value or a branch node in the tree. v is undefined when found is false.
func (t T) getWithTreeNodes(path ...string) (current mapWithValue, found bool) {
	current = t.items
	for i := 0; i < len(path); i++ {
		k := path[i]
		if current.nested == nil {
			return
		}
		var nested mapWithValue
		var ok bool
		if nested, ok = current.nested[k]; !ok {
			found = false
			return
		} else {
			current = nested
		}
	}
	found = true
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
	var previous mapType = nil
	var current mapWithValue
	current = t.items
	i := 0
	for ; i < len(path); i++ {
		k := path[i]
		if current.nested == nil {
			// current branch has no map value initialized, make it
			current.nested = make(mapType)
			// store the "current" object into the map as current is a copy
			if previous == nil {
				t.items = current
			} else {
				previous[path[i-1]] = current
			}
		}
		var next mapWithValue
		var ok bool
		if next, ok = current.nested[k]; !ok {
			// try to traverse into the map at the current path
			// we didn't find one, let's make a blank one
			current.nested[k] = next
		}
		previous = current.nested
		current = next
	}

	// set the value
	if current.value == nil {
		// only add the item count if this element was missing,
		// replacing elements does not increase the count
		t.itemCount++
	}
	// set the current value
	current.value = value

	// store the "current" object into the map as current is a copy
	if previous == nil {
		// no previous, so we're at the root, replace the current root
		t.items = current
	} else {
		// we have a previous, which will point to a valid map, set this current to that key.
		previous[path[i-1]] = current
	}
}

func (t *T) Remove(path ...string) {
	if len(path) == 0 {
		// no path given, do nothing
		return
	}
	lastIndex := len(path) - 1
	parentPath := path[0:lastIndex]
	if parent, ok := t.getWithTreeNodes(parentPath...); !ok {
		// not found, nothing to remove
		return
	} else {
		finalPathPart := path[lastIndex]
		branch, found := parent.nested[finalPathPart]
		if found {
			branch.value = nil
			parent.nested[finalPathPart] = branch
			t.itemCount--
			t.deleteEmptyParentPaths(parentPath...)
		}
	}
}

// deleteEmptyParentPaths goes up the tree from a starting point, deleting any branch nodes that are empty
// this reclaims any map memory that is no longer used to prevent memory leaks
func (t *T) deleteEmptyParentPaths(path ...string) {
	stack := make([]*mapWithValue, 0, len(path))
	stack = append(stack, &t.items)
	for _, key := range path {
		if branch, ok := last(stack).nested[key]; !ok {
			// didn't find the branch... stop
			break
		} else {
			// found the branch, traverse
			stack = append(stack, &branch)
		}
	}

	// parent paths collected
	for i := len(stack) - 1; i > 0; i-- {
		branch := stack[i]
		if len(branch.nested) == 0 && branch.value == nil {
			// this branch is empty, remove it
			delete(stack[i-1].nested, path[i-1])
		}
	}
}

// last returns the last item in the stack, convenience method to keep things cleaner looking
func last(stack []*mapWithValue) *mapWithValue {
	return stack[len(stack)-1]
}

func (t *T) Keys() (paths [][]string) {
	if t.Len() == 0 {
		return
	}
	paths = make([][]string, 0, t.Len())
	keysDFSRecursive(&t.items, []string{}, &paths)
	return
}

// keysDFSRecursive collects all of the keys for all leaf nodes in the map
// this method is not very efficient and could be build to be faster.
func keysDFSRecursive(m *mapWithValue, parentPath []string, output *[][]string) {
	for key, child := range m.nested {
		path := make([]string, len(parentPath)+1)
		for i, k := range parentPath {
			path[i] = k
		}
		path[len(path)-1] = key
		if child.value != nil {
			*output = append(*output, path)
		}
		keysDFSRecursive(&child, path, output)
	}
}
