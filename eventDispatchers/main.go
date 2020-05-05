package eventDispatchers

import "google.golang.org/protobuf/proto"

type IEventDispatcher interface {
	Send(object string, version int32, payload proto.Message)
}
