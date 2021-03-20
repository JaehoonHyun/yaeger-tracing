package tracer

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
)

func NewTracer() (opentracing.Tracer, io.Closer, error) {
	// load config from environment variables
	cfg, _ := jaegercfg.FromEnv()

	// create tracer from config
	return cfg.NewTracer(
		config.Metrics(jprom.New()),
	)
}
