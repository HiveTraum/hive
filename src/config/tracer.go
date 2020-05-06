package config

import (
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
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
	log.Log().Msg("OpenTracing successfully initiated")
	return tracer, closer
}
