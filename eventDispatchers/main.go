package eventDispatchers

import "github.com/golang/protobuf/proto"

type IEventDispatcher interface {
	Send(object string, version int32, payload proto.Message)
}
