package functools

import (
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"strings"
)

func UUIDSliceToByteArraySlice(identifiers []uuid.UUID) [][]byte {
	bytes := make([][]byte, len(identifiers))

	for i, id := range identifiers {
		bytes[i] = id.Bytes()
	}

	return bytes
}

func ByteArraySliceToUUIDSlice(bytes [][]byte) []uuid.UUID {
	identifiers := make([]uuid.UUID, len(bytes))

	for i, b := range bytes {

		uid, err := uuid.FromBytes(b)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}

		identifiers[i] = uid
	}

	return identifiers
}

func UUIDListToStringList(id []uuid.UUID) []string {
	identifiers := make([]string, len(id))

	for i, v := range id {
		identifiers[i] = v.String()
	}

	return identifiers
}

func UUIDListToString(id []uuid.UUID, delimiter string) string {
	return strings.Join(UUIDListToStringList(id), delimiter)
}

func UUIDListToPGArray(id []uuid.UUID) string {
	return StringsToPGArray(UUIDListToStringList(id))
}

func StringsSliceToUUIDSlice(id []string) []uuid.UUID {
	identifiers := make([]uuid.UUID, len(id))

	for i, v := range id {
		identifiers[i] = uuid.FromStringOrNil(v)
	}

	return identifiers
}
