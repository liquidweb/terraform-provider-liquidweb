module github.com/liquidweb/terraform-provider-liquidweb

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0

require (
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/hashicorp/terraform v0.12.2
	github.com/liquidweb/go-lwApi v0.0.0-20190605172801-52a4864d2738
	github.com/liquidweb/liquidweb-go v1.6.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/spf13/viper v1.4.0
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible
)

go 1.13
