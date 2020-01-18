package functools

import (
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
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
