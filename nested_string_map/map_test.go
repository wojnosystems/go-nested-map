package nested_string_map

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestT_At(t *testing.T) {
	cases := map[string]struct {
		input         T
		path          []string
		expected      interface{}
		expectedFound bool
	}{
		"with empty map it returns nothing": {
			input: T{},
			path:  []string{"first", "second"},
		},
		"with item at root it succeeds": {
			input: func() (m T) {
				m = T{}
				m.Put(5, RootPath...)
				return
			}(),
			path:          RootPath,
			expected:      5,
			expectedFound: true,
		},
		"with item with invalid path it returns nothing": {
			input: func() (m T) {
				m = T{}
				m.Put(5, "invalid", "path")
				return
			}(),
			path: []string{"first", "second"},
		},
		"with item with valid path it returns the item": {
			input: func() (m T) {
				m = T{}
				m.Put(5, "first", "second")
				return
			}(),
			path:          []string{"first", "second"},
			expected:      5,
			expectedFound: true,
		},
		"with too few items with valid path it does not return the tree node": {
			input: func() (m T) {
				m = T{}
				m.Put(5, "first", "second")
				return
			}(),
			path:          []string{"first"},
			expectedFound: false,
		},
		"with items with shared parent paths it returns the item": {
			input: func() (m T) {
				m = T{}
				m.Put(2, "first", "second")
				m.Put(3, "first", "third")
				m.Put(4, "first", "fourth")
				return
			}(),
			path:          []string{"first", "second"},
			expected:      2,
			expectedFound: true,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual, ok := c.input.Get(c.path...)
			if c.expectedFound {
				require.True(t, ok)
				assert.Equal(t, c.expected, actual)
			} else {
				assert.False(t, ok)
			}
		})
	}
}

func TestT_Remove(t *testing.T) {
	cases := map[string]struct {
		input T
		path  []string
	}{
		"empty path": {
			input: T{},
			path:  []string{},
		},
		"remove non-existent item": {
			input: T{},
			path:  []string{"first", "second"},
		},
		"remove non-existent item initialized": {
			input: T{},
			path:  []string{"first", "second"},
		},
		"removes item initialized": {
			input: func() (m T) {
				m = T{}
				m.Put(5, "first", "second")
				return
			}(),
			path: []string{"first", "second"},
		},
		"removes item not-initialized": {
			input: func() (m T) {
				m = T{}
				m.Put(5, "first", "second")
				return
			}(),
			path: []string{"first", "second"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.input.Remove(c.path...)
			_, ok := c.input.Get(c.path...)
			assert.False(t, ok)
		})
	}
}

func TestT_Len(t *testing.T) {
	actual := T{}
	assert.Equal(t, 0, actual.Len())
	actual.Put(5, "path", "to", "anything")
	assert.Equal(t, 1, actual.Len())
	actual.Put(6, "path", "to", "another", "thing")
	assert.Equal(t, 2, actual.Len())
	actual.Remove("path", "to", "anything")
	assert.Equal(t, 1, actual.Len())
	actual.Remove("path", "to", "another", "thing")
	assert.Equal(t, 0, actual.Len())
}

func TestT_IsEmpty(t *testing.T) {
	actual := T{}
	assert.True(t, actual.IsEmpty())
	actual.Put(5, "path", "to", "anything")
	assert.False(t, actual.IsEmpty())
	actual.Remove("path", "to", "anything")
	assert.True(t, actual.IsEmpty())
}

func TestT_Keys(t *testing.T) {
	cases := map[string]struct {
		input    T
		expected [][]string
	}{
		"empty": {
			input:    T{},
			expected: [][]string{},
		},
		"with uninitialized empty": {
			input:    T{},
			expected: [][]string{},
		},
		"with an item": {
			input: func() (m T) {
				m.Put(5, "a", "b", "c", "d")
				return
			}(),
			expected: [][]string{
				{"a", "b", "c", "d"},
			},
		},
		"with many items": {
			input: func() (m T) {
				m.Put(5, "a", "b", "c", "d")
				m.Put(5, "a", "b", "c", "e")
				m.Put(5, "a", "b", "c", "f")
				m.Put(5, "a", "b", "g", "d")
				m.Put(5, "a", "b", "h", "d")
				return
			}(),
			expected: [][]string{
				{"a", "b", "c", "d"},
				{"a", "b", "c", "e"},
				{"a", "b", "c", "f"},
				{"a", "b", "g", "d"},
				{"a", "b", "h", "d"},
			},
		},
		"with some items removed": {
			input: func() (m T) {
				m.Put(5, "a", "b", "c", "d")
				m.Put(5, "a", "b", "c", "e")
				m.Put(5, "remove", "b", "c", "f")
				m.Put(5, "a", "b", "g", "d")
				m.Put(5, "remove", "b", "h", "d")
				m.Remove("remove", "b", "c", "f")
				m.Remove("remove", "b", "h", "d")
				return
			}(),
			expected: [][]string{
				{"a", "b", "c", "d"},
				{"a", "b", "c", "e"},
				{"a", "b", "g", "d"},
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.input.Keys()
			for _, key := range c.expected {
				assert.Contains(t, actual, key)
			}
			if len(c.expected) == 0 {
				assert.Empty(t, actual)
			}
		})
	}
}

func Test_deleteEmptyParentPathsRejectsMissingBranches(t *testing.T) {
	actual := T{}
	actual.Put(5, "some", "path")
	actual.deleteEmptyParentPaths("some", "missing-path")
	_, ok := actual.Get("some", "path")
	assert.True(t, ok)
}

func Test_deleteEmptyParentPathsRemoveEmptyBranches(t *testing.T) {
	actual := T{}
	actual.Put(nil, "some", "path", "location")
	actual.deleteEmptyParentPaths("some", "path", "location")
	_, ok := actual.Get("some", "path", "location")
	assert.False(t, ok)
}

func Test_deleteEmptyParentPathsIgnoresUninitializedItems(t *testing.T) {
	assert.NotPanics(t, func() {
		(&T{}).deleteEmptyParentPaths("some", "missing-path")
	})
}
