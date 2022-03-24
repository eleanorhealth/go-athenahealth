package tracer

import (
	"net/http"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

type config struct {
	service string
}

type Option func(*config)

func WithService(service string) Option {
	return func(cfg *config) {
		cfg.service = service
	}
}

func WrapRoundTripper(transport http.RoundTripper, opts ...Option) http.RoundTripper {
	cfg := &config{
		service: "go-athenahealth.client",
	}
	for _, opt := range opts {
		opt(cfg)
	}

	rtOpts := []httptrace.RoundTripperOption{
		httptrace.WithBefore(func(req *http.Request, span ddtrace.Span) {
			span.SetTag(ext.ServiceName, cfg.service)
			span.SetTag(ext.SpanType, ext.SpanTypeHTTP)
			span.SetTag(ext.HTTPMethod, req.Method)
			span.SetTag(ext.HTTPURL, req.URL.Path)
		}),
	}

	return httptrace.WrapRoundTripper(transport, rtOpts...)
}
