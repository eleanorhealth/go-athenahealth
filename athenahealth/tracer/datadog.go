package tracer

import (
	"context"
	"net/url"
	"regexp"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var idRegex = regexp.MustCompile(`(/)(\d+)(/?)`)

type Datadog struct {
	service string
}

type Option func(*Datadog)

func NewDatadog(opts ...Option) *Datadog {
	d := &Datadog{
		service: "go-athenahealth.client",
	}

	for _, o := range opts {
		o(d)
	}

	return d
}

func WithService(service string) Option {
	return func(d *Datadog) {
		d.service = service
	}
}

func (d *Datadog) Before(ctx context.Context, method, path string) context.Context {
	path = cleanPath(path)

	opts := []ddtrace.StartSpanOption{
		tracer.SpanType(ext.SpanTypeHTTP),
		tracer.Tag(ext.HTTPMethod, method),
		tracer.Tag(ext.HTTPURL, path),
		tracer.ServiceName(d.service),
	}

	_, ctx = tracer.StartSpanFromContext(ctx, "go-athenahealth", opts...)

	return ctx
}

func (d *Datadog) After(ctx context.Context) {
	if span, ok := tracer.SpanFromContext(ctx); ok {
		span.Finish()
	}
}

func cleanPath(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		return ""
	}

	return idRegex.ReplaceAllString(u.Path, "$1:id:$3")
}
