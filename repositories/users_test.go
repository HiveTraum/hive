package repositories

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetRawQueryWithModifiedLimit(t *testing.T) {
	t.Parallel()

	repeatedUserID := uuid.NewV4()

	q := GetUsersQuery{
		Limit: 100,
		Id:    []uuid.UUID{uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), repeatedUserID, repeatedUserID},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Limit == 6)
}

func TestGetRawQueryLimitWithEmptyId(t *testing.T) {
	t.Parallel()

	q := GetUsersQuery{
		Limit: 100,
		Id:    []uuid.UUID{},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Limit == 100)
}

func TestGetRawQueryWithLimitLessThenId(t *testing.T) {
	t.Parallel()

	repeatedUserID := uuid.NewV4()

	q := GetUsersQuery{
		Limit: 3,
		Id:    []uuid.UUID{uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), repeatedUserID, repeatedUserID},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Limit == 3)
}

func TestGetRawQueryWithEmptyId(t *testing.T) {
	t.Parallel()

	q := GetUsersQuery{
		Id: []uuid.UUID{},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Id == "{}")
}

func TestGetRawQueryWithId(t *testing.T) {
	t.Parallel()

	first, second, third := uuid.NewV4(), uuid.NewV4(), uuid.NewV4()

	q := GetUsersQuery{
		Id: []uuid.UUID{first, second, third},
	}

	rw := convertGetUsersQueryToRaw(q)

	require.True(t, rw.Id == fmt.Sprintf("{%s,%s,%s}", first.String(), second.String(), third.String()))
}
