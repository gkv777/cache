package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func get2Q(t *testing.T, n1, n2, n3 int) (*TwoQCache, error) {
	lru, err := newLru(n3)
	if err != nil {
		return nil, err
	}
	fifo, err := newFifo(n1, n2)
	if err != nil {
		return nil, err
	}
	return &TwoQCache{
		fifo: fifo,
		lru:  lru,
	}, nil
}

func TestNew2QCache(t *testing.T) {
	tcs := map[string]struct {
		n       int
		err     error
		lenLru  int
		lenFifo int
	}{
		"ok_10":   {10, nil, 2, 8},
		"err_0":   {0, ErrCapSize, 0, 0},
		"err_low": {4, ErrCapSize, 0, 0},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			cache, err := NewTwoQCache(tc.n)
			if tc.err != nil {
				require.ErrorIs(t, tc.err, err)
				require.Nil(t, cache)
			} else {
				require.NoError(t, err)
				//require.Equal(t, tc.lenLru, cache.lru.Len())
			}
		})
	}
}

func TestTwoQCacheAdd(t *testing.T) {
	cache, _ := get2Q(t, 2, 6, 2)
	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("key%d", i)
		ok := cache.Add(key, key)
		t.Run(key, func(t *testing.T) {
			require.Equal(t, true, ok)
			if i <= 2 {
				require.Equal(t, i, cache.fifo.inLen())
				require.Equal(t, 0, cache.fifo.outLen())
				require.Equal(t, 0, cache.lru.Len())
			} else if i <= 8 {
				require.Equal(t, 2, cache.fifo.inLen())
				require.Equal(t, i-2, cache.fifo.outLen())
				require.Equal(t, 0, cache.lru.Len())
			} else {
				require.Equal(t, 2, cache.fifo.inLen())
				require.Equal(t, 6, cache.fifo.outLen())
				require.Equal(t, 0, cache.lru.Len())
			}
		})

	}
	t.Run("already exists", func(t *testing.T) {
		ok := cache.Add("key7", "key7")
		require.Equal(t, false, ok)
	})

	t.Run("key5 move to Am and check", func(t *testing.T) {
		cache.Get("key5")
		ok := cache.Add("key5", "key5")
		require.Equal(t, false, ok)
	})
}

func TestTwoQCacheGet(t *testing.T) {
	tcs := []struct {
		name string
		id   int
		ok   bool
		in   int
		out  int
		am   int
	}{
		{"key8 in A1in", 8, true, 2, 6, 0},
		{"key6 in A1out", 6, true, 2, 5, 1},
		{"key5 in A1out", 5, true, 2, 4, 2},
		{"key4 in A1out", 4, true, 2, 3, 2},
		{"key6 not in Am", 6, false, 2, 3, 2},
		{"key4 in Am", 4, true, 2, 3, 2},
	}

	cache, _ := get2Q(t, 2, 6, 2)
	for i := 1; i <= 8; i++ {
		key := fmt.Sprintf("key%d", i)
		_ = cache.Add(key, key)
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := cache.Get(fmt.Sprintf("key%d", tc.id))
			require.Equal(t, tc.ok, ok)
			if tc.ok {
				require.Equal(t, fmt.Sprintf("key%d", tc.id), val)
			}
			require.Equal(t, tc.in, cache.fifo.inLen())
			require.Equal(t, tc.out, cache.fifo.outLen())
			require.Equal(t, tc.am, cache.lru.Len())
		})
	}

}

func TestTwoQCacheRemove(t *testing.T) {
	cache, _ := get2Q(t, 2, 6, 2)
	for i := 1; i <= 8; i++ {
		key := fmt.Sprintf("key%d", i)
		_ = cache.Add(key, key)
	}

	t.Run("key8 remove from A1in", func(t *testing.T) {
		ok := cache.Remove("key8")
		require.Equal(t, true, ok)
		require.Equal(t, 1, cache.fifo.inLen())
		require.Equal(t, 6, cache.fifo.outLen())
		require.Equal(t, 0, cache.lru.Len())
	})

	t.Run("key8 remove from A1out", func(t *testing.T) {
		ok := cache.Remove("key5")
		require.Equal(t, true, ok)
		require.Equal(t, 1, cache.fifo.inLen())
		require.Equal(t, 5, cache.fifo.outLen())
		require.Equal(t, 0, cache.lru.Len())
	})

	cache.Get("key4")
	cache.Get("key3")

	t.Run("key3 remove from Am", func(t *testing.T) {
		ok := cache.Remove("key3")
		require.Equal(t, true, ok)
		require.Equal(t, 1, cache.fifo.inLen())
		require.Equal(t, 3, cache.fifo.outLen())
		require.Equal(t, 1, cache.lru.Len())
	})

	t.Run("key3 not exists", func(t *testing.T) {
		ok := cache.Remove("key3")
		require.Equal(t, false, ok)
		require.Equal(t, 1, cache.fifo.inLen())
		require.Equal(t, 3, cache.fifo.outLen())
		require.Equal(t, 1, cache.lru.Len())
	})
}
