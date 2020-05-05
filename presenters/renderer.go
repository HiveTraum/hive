package presenters

import (
	"hive/enums"
	"hive/repositories"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
	"reflect"
)

type Writer func(proto.Message) ([]byte, error)

type Renderer struct {
	contentTypeParsers map[enums.ContentType]Writer
}

func (renderer *Renderer) render(w http.ResponseWriter, r *http.Request, status int, data proto.Message) error {

	contentType := repositories.GetContentTypeHeader(r)

	contentTypeRenderer := renderer.contentTypeParsers[contentType]
	if contentTypeRenderer == nil {
		return errors.New(fmt.Sprintf("renderer not found for %s", contentType))
	}

	if data != nil && !reflect.ValueOf(data).IsNil() {
		dataBytes, err := contentTypeRenderer(data)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		} else {
			w.WriteHeader(status)
			_, _ = w.Write(dataBytes)
		}
	} else {
		w.WriteHeader(status)
	}

	return nil
}

func (renderer *Renderer) Render(w http.ResponseWriter, r *http.Request, status int, data proto.Message) {
	err := renderer.render(w, r, status, data)
	if err != nil {
		sentry.CaptureException(err)
	}
}

func InitRenderer() *Renderer {

	return &Renderer{contentTypeParsers: map[enums.ContentType]Writer{
		enums.JSONContentType:   protojson.Marshal,
		enums.BinaryContentType: proto.Marshal,
	}}
}
