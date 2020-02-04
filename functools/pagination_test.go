package functools

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLimitPageToLimitOffset(t *testing.T) {
	t.Parallel()
	limit, offset := LimitPageToLimitOffset(1, 1)
	require.Equal(t, 1, limit)
	require.Equal(t, 0, offset)
}

func TestLimitPageToLimitOffsetWithBiggerLimit(t *testing.T) {
	t.Parallel()
	limit, offset := LimitPageToLimitOffset(5, 1)
	require.Equal(t, 5, limit)
	require.Equal(t, 0, offset)
}

func TestLimitPageToLimitOffsetWithBiggerPage(t *testing.T) {
	t.Parallel()
	limit, offset := LimitPageToLimitOffset(1, 5)
	require.Equal(t, 1, limit)
	require.Equal(t, 4, offset)
}

func TestLimitPageToLimitOffsetWithBiggerPageAndLimit(t *testing.T) {
	t.Parallel()
	limit, offset := LimitPageToLimitOffset(5, 5)
	require.Equal(t, 5, limit)
	require.Equal(t, 20, offset)
}

func TestLimitPageToLimitOffsetWithZeroPage(t *testing.T) {
	t.Parallel()
	limit, offset := LimitPageToLimitOffset(5, 0)
	require.Equal(t, 5, limit)
	require.Equal(t, 0, offset)
}

func TestLimitPageToLimitOffsetWithZeroLimit(t *testing.T) {
	t.Parallel()
	limit, offset := LimitPageToLimitOffset(0, 1)
	require.Equal(t, 1, limit)
	require.Equal(t, 0, offset)
}

func TestHasNext(t *testing.T) {
	t.Parallel()
	require.True(t, HasNext(6, 2, 1))
	require.True(t, HasNext(6, 2, 2))
	require.False(t, HasNext(6, 2, 3))
	require.True(t, HasNext(5, 2, 2))
	require.False(t, HasNext(5, 3, 2))
	require.True(t, HasNext(6, 0, 1))
	require.True(t, HasNext(6, 0, 5))
	require.False(t, HasNext(6, 0, 6))
	require.False(t, HasNext(0, 1, 1))
	require.False(t, HasNext(0, 1, 1))
}

func TestHasPrevious(t *testing.T) {
	t.Parallel()
	require.False(t, HasPrevious(1))
	require.True(t, HasPrevious(2))
}
