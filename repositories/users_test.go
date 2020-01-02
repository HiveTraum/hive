package repositories

import (
	"auth/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetRawQueryWithModifiedLimit(t *testing.T) {
	t.Parallel()

	q := GetUsersQuery{
		Limit:       100,
		Id: []models.UserID{1, 2, 3, 4, 5, 5},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Limit == 6)
}

func TestGetRawQueryLimitWithEmptyId(t *testing.T) {
	t.Parallel()

	q := GetUsersQuery{
		Limit:       100,
		Id: []models.UserID{},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Limit == 100)
}

func TestGetRawQueryWithLimitLessThenId(t *testing.T) {
	t.Parallel()

	q := GetUsersQuery{
		Limit:       3,
		Id: []models.UserID{1, 2, 3, 4, 5, 5},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Limit == 3)
}

func TestGetRawQueryWithEmptyId(t *testing.T) {
	t.Parallel()

	q := GetUsersQuery{
		Id: []models.UserID{},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Id == "{}")
}

func TestGetRawQueryWithId(t *testing.T) {
	t.Parallel()

	q := GetUsersQuery{
		Id: []models.UserID{1, 2, 3, 4, 5},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Id == "{1,2,3,4,5}")
}
