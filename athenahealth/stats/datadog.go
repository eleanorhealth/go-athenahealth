package stats

import (
	"net/url"
	"regexp"

	"github.com/DataDog/datadog-go/statsd"
)

var idRegex = regexp.MustCompile(`(/)(\d+)(/?)`)

type Datadog struct {
	client statsd.ClientInterface
}

func NewDatadog(client statsd.ClientInterface) *Datadog {
	return &Datadog{
		client: client,
	}
}

func (d *Datadog) Request(method, path string) error {
	path = cleanPath(path)

	return d.client.Incr("athenahealth.requests", []string{
		"http_method:" + method,
		"http_path:" + path,
	}, 1.0)
}

func (d *Datadog) ResponseSuccess() error {
	return d.client.Incr("athenahealth.responses.success", []string{}, 1.0)
}

func (d *Datadog) ResponseError() error {
	return d.client.Incr("athenahealth.responses.error", []string{}, 1.0)
}

func cleanPath(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		return ""
	}

	return idRegex.ReplaceAllString(u.Path, "$1:id:$3")
}
