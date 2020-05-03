package eventDispatchers

import (
	"auth/config"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
)

type NSQEventDispatcher struct {
	producer    *nsq.Producer
	environment *config.Environment
}

func InitNSQEventDispatcher(producer *nsq.Producer, environment *config.Environment) *NSQEventDispatcher {
	return &NSQEventDispatcher{producer: producer, environment: environment}
}

func (dispatcher NSQEventDispatcher) Send(object string, version int32, payload proto.Message) {
	topic := fmt.Sprintf("%s-%s-%d", dispatcher.environment.ESBSender, object, version)

	data, err := proto.Marshal(payload)

	if err != nil {
		sentry.CaptureException(err)
		return
	}

	err = dispatcher.producer.Publish(topic, data)

	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
