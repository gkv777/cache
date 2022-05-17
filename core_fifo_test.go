package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFifo(t *testing.T) {
	f, err := newFifo(2, 4)
	require.NoError(t, err)
	require.Equal(t, 6, f.Cap())
	require.Equal(t, 0, f.Len())
	f, err = newFifo(0, 1)
	require.Error(t, err)
	require.Nil(t, f)
}

func Test_fifoAdd(t *testing.T) {
	l, _ := newFifo(2, 2)

	t.Run("key1", func(t *testing.T) {
		ok := l.Add("key1", "value1")
		require.Equal(t, true, ok)
		require.Equal(t, 1, l.Len())
		require.Equal(t, 1, l.inLen())
		require.Equal(t, 0, l.outLen())
	})

	t.Run("key2", func(t *testing.T) {
		ok := l.Add("key2", "value2")
		require.Equal(t, true, ok)
		require.Equal(t, 2, l.Len())
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 0, l.outLen())
	})

	t.Run("key2_duplicate", func(t *testing.T) {
		ok := l.Add("key2", "value2")
		require.Equal(t, false, ok)
		require.Equal(t, 2, l.Len())
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 0, l.outLen())
	})

	t.Run("key3", func(t *testing.T) {
		ok := l.Add("key3", "value3")
		require.Equal(t, true, ok)
		require.Equal(t, 3, l.Len())
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 1, l.outLen())
	})

	t.Run("key4", func(t *testing.T) {
		ok := l.Add("key4", "value4")
		require.Equal(t, true, ok)
		require.Equal(t, 4, l.Len())
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 2, l.outLen())
	})

	t.Run("key5", func(t *testing.T) {
		ok := l.Add("key5", "value5")
		require.Equal(t, true, ok)
		require.Equal(t, 4, l.Len())
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 2, l.outLen())
	})
}

func Test_fifoGet(t *testing.T) {
	items := []struct {
		key string
		val string
	}{{"key1", "value1"}, {"key2", "value2"}, {"key3", "value3"}, {"key4", "value4"}}

	l, _ := newFifo(2,3)
	for _, i := range items {
		l.Add(i.key, i.val)
	}

	// in A1in
	t.Run("key3", func(t *testing.T) {
		val, ok, out := l.Get("key3")
		require.Equal(t, true, ok)
		require.Equal(t, false, out)
		require.Equal(t, "value3", val)
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 2, l.outLen())
	})

	// in A1out
	t.Run("key1", func(t *testing.T) {
		val, ok, out := l.Get("key1")
		require.Equal(t, true, ok)
		require.Equal(t, true, out)
		require.Equal(t, "value1", val)
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 1, l.outLen())
	})

	// not in cache
	t.Run("key1", func(t *testing.T) {
		val, ok, out := l.Get("key1")
		require.Equal(t, false, ok)
		require.Equal(t, false, out)
		require.Equal(t, "", val)
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 1, l.outLen())
	})
	
}

func Test_fifoRemove(t *testing.T) {
	items := []struct {
		key string
		val string
	}{{"key1", "value1"}, {"key2", "value2"}, {"key3", "value3"}, {"key4", "value4"}}

	l, _ := newFifo(2,2)
	for _, i := range items {
		l.Add(i.key, i.val)
	}

	t.Run("key1_exists", func(t *testing.T) {
		ok := l.Remove("key1")
		require.Equal(t, true, ok)
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 1, l.outLen())
	})

	t.Run("key1_absence", func(t *testing.T) {
		ok := l.Remove("key1")
		require.Equal(t, false, ok)
		require.Equal(t, 2, l.inLen())
		require.Equal(t, 1, l.outLen())
	})

	t.Run("key3", func(t *testing.T) {
		ok := l.Remove("key3")
		require.Equal(t, true, ok)
		require.Equal(t, 1, l.inLen())
		require.Equal(t, 1, l.outLen())
	})

}
