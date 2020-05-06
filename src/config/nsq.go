package config

import (
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog/log"
)

func InitNSQProducer(environment *Environment) *nsq.Producer {
	producer, err := nsq.NewProducer(environment.NSQLookupAddress, nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	log.Log().Msg("NSQ successfully initiated")
	return producer
}
