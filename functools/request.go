package functools

import (
	"auth/config"
	"auth/models"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Request struct {
	*http.Request
	Response http.ResponseWriter
}

const (
	JSONContentType   = "application/json"
	BinaryContentType = "application/octet-stream"
)

func (request *Request) GetContentType() string {
	contentType := request.Header.Get("Content-Type")

	if contentType == "" {
		contentType = BinaryContentType
	}

	return contentType
}

func (request *Request) IsContentTypeAllowed(allowedTypes *[]string) bool {
	if allowedTypes == nil {
		allowedTypes = &[]string{BinaryContentType, JSONContentType}
	}

	contentType := request.GetContentType()
	return Contains(contentType, *allowedTypes)
}

func (request *Request) IsProto() bool {
	return request.GetContentType() == BinaryContentType
}

func (request *Request) ParseBody(message proto.Message) error {

	b, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	if request.IsProto() {
		err = proto.Unmarshal(b, message)
	} else {
		err = json.Unmarshal(b, message)
	}

	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	return err
}

func GetLimit(values url.Values) int {
	limitQuery := values.Get("limit")
	if limitQuery == "" {
		env := config.GetEnvironment()
		return env.DefaultPaginationLimit
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		env := config.GetEnvironment()
		return env.DefaultPaginationLimit
	}

	return limit
}

func GetPage(values url.Values) int {
	pageQuery := values.Get("page")
	if pageQuery == "" {
		return 1
	}

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		return 1
	}

	return page
}

func GetPagination(values url.Values) *models.PaginationRequest {
	return &models.PaginationRequest{
		Page:  GetPage(values),
		Limit: GetLimit(values),
	}
}

func (request *Request) GetAuthorizationHeader() string {
	return request.Header.Get("Authorization")
}
