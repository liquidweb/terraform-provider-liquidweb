package liquidweb

import (
	opentracing "github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

func traceError(span opentracing.Span, err error) {
	span.SetTag("error", "true")
	span.LogFields(opentracinglog.String("error", err.Error()))
}
