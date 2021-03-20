package tracer

import (
	"io"

	"github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func NewTracer() (opentracing.Tracer, io.Closer, error) {
	// load config from environment variables
	// cfg := jaegercfg.Configuration{
	// 	ServiceName: serviceName,
	// 	// sampling 정책을 결정
	// 	Sampler: &jaegercfg.SamplerConfig{
	// 		Type:  jaeger.SamplerTypeConst,
	// 		Param: 1,
	// 	},
	// 	Reporter: &jaegercfg.ReporterConfig{
	// 		LogSpans: true,
	// 	},
	// }
	cfg, _ := jaegercfg.FromEnv()

	// create tracer from config
	return cfg.NewTracer(
		jaegercfg.Logger(jaeger.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}
