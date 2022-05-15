package core

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/require"
)

func newLru(t *testing.T) *LRUBasic {
	return &LRUBasic{
		cap:   3,
		items: make(map[string]*list.Element),
		queue: list.New(),
	}
}

func TestNewLRUCache(t *testing.T) {
	l, err := NewLRUBasicCache(4)
	require.NoError(t, err)
	require.Equal(t, 4, l.Cap())
	require.Equal(t, 0, l.Len())
	l, err = NewLRUBasicCache(0)
	require.Error(t, err)
	require.Equal(t, nil, l)
}

func Test_lruAdd(t *testing.T) {
	l := newLru(t)

	t.Run("key1", func(t *testing.T) {
		ok := l.Add("key1", "value1")
		require.Equal(t, true, ok)
		require.Equal(t, 1, l.Len())
		require.Equal(t, l.getFirst().key, "key1")
		require.Equal(t, l.getLast().key, "key1")
	})

	t.Run("key2", func(t *testing.T) {
		ok := l.Add("key2", "value2")
		require.Equal(t, true, ok)
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.getFirst().key, "key2")
		require.Equal(t, l.getLast().key, "key1")
	})

	t.Run("key2_duplicate", func(t *testing.T) {
		ok := l.Add("key2", "value2")
		require.Equal(t, false, ok)
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.getFirst().key, "key2")
		require.Equal(t, l.getLast().key, "key1")
	})

	t.Run("key3", func(t *testing.T) {
		ok := l.Add("key3", "value3")
		require.Equal(t, true, ok)
		require.Equal(t, 3, l.Len())
		require.Equal(t, l.getFirst().key, "key3")
		require.Equal(t, l.getLast().key, "key1")
	})

	t.Run("key4", func(t *testing.T) {
		ok := l.Add("key4", "value4")
		require.Equal(t, true, ok)
		require.Equal(t, 3, l.Len())
		require.Equal(t, l.getFirst().key, "key4")
		require.Equal(t, l.getLast().key, "key2")
	})
}

func Test_lruGet(t *testing.T) {
	items := []struct {
		key string
		val string
	}{{"key1", "value1"},{"key2", "value2"},{"key3", "value3"},{"key4", "value4"},}

	l := newLru(t)
	for _, i := range items {
		l.Add(i.key, i.val)
	}

	t.Run("key1", func(t *testing.T) {
		val, ok := l.Get("key1")
		require.Equal(t, false, ok)
		require.Equal(t, "", val)
		require.Equal(t, "key4", l.getFirst().key)
		require.Equal(t, "key2", l.getLast().key)
	})

	t.Run("key3", func(t *testing.T) {
		val, ok := l.Get("key3")
		require.Equal(t, true, ok)
		require.Equal(t, "value3", val)
		require.Equal(t, "key3", l.getFirst().key)
		require.Equal(t, "key2", l.getLast().key)
	})
}

func Test_lruRemove(t *testing.T) {
	items := []struct {
		key string
		val string
	}{{"key1", "value1"},{"key2", "value2"},{"key3", "value3"},{"key4", "value4"},}

	l := newLru(t)
	for _, i := range items {
		l.Add(i.key, i.val)
	}

	t.Run("key1", func(t *testing.T) {
		ok := l.Remove("key1")
		require.Equal(t, false, ok)
		require.Equal(t, 3, l.Len())
		require.Equal(t, "key4", l.getFirst().key)
		require.Equal(t, "key2", l.getLast().key)
	})

	t.Run("key3", func(t *testing.T) {
		ok := l.Remove("key3")
		require.Equal(t, true, ok)
		require.Equal(t, 2, l.Len())
		require.Equal(t, "key4", l.getFirst().key)
		require.Equal(t, "key2", l.getLast().key)
	})

	t.Run("key3_after", func(t *testing.T) {
		ok := l.Remove("key3")
		require.Equal(t, false, ok)
		require.Equal(t, 2, l.Len())
		require.Equal(t, "key4", l.getFirst().key)
		require.Equal(t, "key2", l.getLast().key)
	})
}
