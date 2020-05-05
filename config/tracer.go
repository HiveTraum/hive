package config

import (
	"github.com/opentracing/opentracing-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"io"
)

func InitTracing(environment *Environment) (opentracing.Tracer, io.Closer) {
	tracer, closer, err := jaegerConfig.Configuration{
		ServiceName: environment.ServiceName,
		RPCMetrics:  true,
	}.NewTracer()

	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
