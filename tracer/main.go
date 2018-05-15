package tracer

import (
	"io"
	"log"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

// New returns a new tracer
func New(serviceName, hostPort string) (opentracing.Tracer, io.Closer) {

	zipkinPropagator := jaeger.NewCombinedB3HTTPHeaderPropagator()

	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  hostPort, // localhost:5775
		},
	}
	tracer, closer, err := cfg.New(
		serviceName,
		config.Logger(jaeger.StdLogger),
		config.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		config.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		config.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	j_tracer, j_closer, j_err := cfg.New(
		serviceName,
		config.Logger(jaeger.StdLogger),
	)
	if j_err != nil {
		log.Fatal(j_err)
	}

	if true {
		return tracer, closer
	} else {
		return j_tracer, j_closer
	}
}
