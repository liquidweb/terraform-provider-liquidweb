package main

import (
	"log"
	"os"
	"github.com/liquidweb/terraform-provider-liquidweb/liquidweb"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	metrics "github.com/uber/jaeger-lib/metrics"
)

func main() {
	log.Printf("Tracing enabled: %t", os.Getenv("JAEGER_DISABLED") == "false")

	cfg, err := jaegercfg.FromEnv()
	cfg.ServiceName = "terraform-provider-liquidweb"
	cfg.Sampler = &jaegercfg.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 1,
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return liquidweb.Provider()
		},
	})
}
