package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew2QCache(t *testing.T) {
	tcs := map[string]struct {
		n int
		err error
		lenLru int
		lenFifo int
	}{
		"ok_10": {10, nil, 2, 8},
		"err_0": {0, ErrCapSize, 0, 0},
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