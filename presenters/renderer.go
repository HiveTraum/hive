package presenters

import (
	"auth/enums"
	"auth/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
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
		bytes, err := contentTypeRenderer(data)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		} else {
			w.WriteHeader(status)
			_, _ = w.Write(bytes)
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
		enums.JSONContentType: func(message proto.Message) ([]byte, error) {
			return json.Marshal(message)
		},
		enums.BinaryContentType: proto.Marshal,
	}}
}
