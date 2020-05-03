package presenters

import (
	"auth/enums"
	"auth/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
)

type Reader func([]byte, proto.Message) error

type Parser struct {
	contentTypeParsers map[enums.ContentType]Reader
}

func (p *Parser) parse(r *http.Request, message proto.Message) error {
	contentType := repositories.GetContentTypeHeader(r)
	contentTypeParser := p.contentTypeParsers[contentType]

	if contentTypeParser == nil {
		return errors.New(fmt.Sprintf("parser not found for %s", contentType))
	}

	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	return contentTypeParser(bytes, message)
}

func (p *Parser) Parse(r *http.Request, w http.ResponseWriter, message proto.Message) error {
	err := p.parse(r, message)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sentry.CaptureException(err)
		return err
	}

	return nil
}

func InitParser() *Parser {
	return &Parser{contentTypeParsers: map[enums.ContentType]Reader{
		enums.JSONContentType: func(bytes []byte, message proto.Message) error {
			return json.Unmarshal(bytes, message)
		},
		enums.BinaryContentType: proto.Unmarshal,
	}}
}
