package events

import (
	"auth/inout"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"net/http"
)

type EventDispatcher struct {
	Url string
}

func (dispatcher *EventDispatcher) Send(event inout.Event) {

	body, err := json.Marshal(event)

	if err != nil {
		sentry.CaptureException(err)
		return
	}

	if dispatcher.Url == "" {
		fmt.Println(string(body))
	} else {
		resp, err := http.Post(dispatcher.Url, "application/json", bytes.NewBuffer(body))

		if err != nil {
			sentry.CaptureException(err)
		}

		resp.Body.Close()
	}
}
