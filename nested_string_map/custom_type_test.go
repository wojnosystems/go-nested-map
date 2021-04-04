package nested_string_map

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type custom struct {
	T
}

func (c *custom) GetInt(path ...string) (value int, ok bool) {
	v, ok := c.Get(path...)
	if ok {
		value, ok = v.(int)
	}
	return
}

func TestCustomType(t *testing.T) {
	cases := map[string]struct {
		path          []string
		input         custom
		expected      int
		expectedFound bool
	}{
		"with empty map it returns nothing": {
			input: custom{},
			path:  []string{"first", "second"},
		},
		"with item with invalid path it returns nothing": {
			path: []string{"first", "second"},
			input: func() (m custom) {
				m = custom{}
				m.Put(5, "invalid", "path")
				return
			}(),
		},
		"with item with valid path it returns the item": {
			path: []string{"first", "second"},
			input: func() (m custom) {
				m = custom{}
				m.Put(5, "first", "second")
				return
			}(),
			expected:      5,
			expectedFound: true,
		},
		"with too few items with valid path it returns the tree node": {
			path: []string{"first"},
			input: func() (m custom) {
				m = custom{}
				m.Put(5, "first", "second")
				return
			}(),
			expectedFound: false,
		},
		"with items with shared parent paths it returns the item": {
			path: []string{"first", "second"},
			input: func() (m custom) {
				m = custom{}
				m.Put(2, "first", "second")
				m.Put(3, "first", "third")
				m.Put(4, "first", "fourth")
				return
			}(),
			expected:      2,
			expectedFound: true,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual, ok := c.input.GetInt(c.path...)
			if c.expectedFound {
				require.True(t, ok)
				assert.Equal(t, c.expected, actual)
			} else {
				assert.False(t, ok)
			}
		})
	}
}
