package stats

import (
	"regexp"
	"strings"

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
	path = removeIDsFromPath(path)

	return d.client.Incr("athenahealth.requests", []string{
		"path:" + strings.ToUpper(method) + " " + strings.ToLower(path),
	}, 1.0)
}

func (d *Datadog) ResponseSuccess() error {
	return d.client.Incr("athenahealth.responses.success", []string{}, 1.0)
}

func (d *Datadog) ResponseError() error {
	return d.client.Incr("athenahealth.responses.error", []string{}, 1.0)
}

func removeIDsFromPath(path string) string {
	return idRegex.ReplaceAllString(path, "$1{id}$3")
}
