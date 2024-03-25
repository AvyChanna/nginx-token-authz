package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	s := New[string]()

	s.Add("a")

	require.Equal(t, 1, s.Len())
	assert.True(t, s.Contains("a"))
	assert.False(t, s.Contains("b"))
	assert.Equal(t, []string{"a"}, s.Values())

	s.Add("a")
	require.Equal(t, 1, s.Len())

	s.Add("b")
	require.Equal(t, 2, s.Len())
	assert.True(t, s.Contains("b"))
	assert.ElementsMatch(t, []string{"a", "b"}, s.Values())

	s.Remove("a")
	require.Equal(t, 1, s.Len())

	s.Remove("b")
	require.Equal(t, 0, s.Len())
}
