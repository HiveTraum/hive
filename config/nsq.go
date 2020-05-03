package config

import "github.com/nsqio/go-nsq"

func InitNSQProducer(environment *Environment) *nsq.Producer {
	producer, err := nsq.NewProducer(environment.NSQLookupAddress, nsq.NewConfig())

	if err != nil {
		panic(err)
	}

	return producer
}
