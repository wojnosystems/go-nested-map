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
		"with item with invalid path it returns nothing": {
			input: func() (m T) {
				m = T{}
				m.Set(5, "invalid", "path")
				return
			}(),
			path: []string{"first", "second"},
		},
		"with item with valid path it returns the item": {
			input: func() (m T) {
				m = T{}
				m.Set(5, "first", "second")
				return
			}(),
			path:          []string{"first", "second"},
			expected:      5,
			expectedFound: true,
		},
		"with too few items with valid path it returns the tree node": {
			input: func() (m T) {
				m = T{}
				m.Set(5, "first", "second")
				return
			}(),
			path: []string{"first"},
			expected: func() (m mapType) {
				m = make(mapType)
				m["second"] = 5
				return
			}(),
			expectedFound: true,
		},
		"with items with shared parent paths it returns the item": {
			input: func() (m T) {
				m = T{}
				m.Set(2, "first", "second")
				m.Set(3, "first", "third")
				m.Set(4, "first", "fourth")
				return
			}(),
			path:          []string{"first", "second"},
			expected:      2,
			expectedFound: true,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual, ok := c.input.At(c.path...)
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
			input: T{items: make(mapType)},
			path:  []string{"first", "second"},
		},
		"removes item initialized": {
			input: func() (m T) {
				m = T{items: make(mapType)}
				m.Set(5, "first", "second")
				return
			}(),
			path: []string{"first", "second"},
		},
		"removes item not-initialized": {
			input: func() (m T) {
				m = T{}
				m.Set(5, "first", "second")
				return
			}(),
			path: []string{"first", "second"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.input.Remove(c.path...)
			_, ok := c.input.At(c.path...)
			assert.False(t, ok)
		})
	}
}
