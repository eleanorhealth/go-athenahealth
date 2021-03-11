package stats

import (
	"github.com/DataDog/datadog-go/statsd"
)

const datadogDefaultNamespace = "athenahealth-client"

type Datadog struct {
	client statsd.ClientInterface
}

func NewDatadog(client statsd.ClientInterface) *Datadog {
	return &Datadog{
		client: client,
	}
}

func (d *Datadog) Request() error {
	err := d.client.Count("athenahealth.requests.count", 1, []string{}, 1.0)
	if err != nil {
		return err
	}

	return d.client.Incr("athenahealth.requests", []string{}, 1.0)
}

func (d *Datadog) ResponseSuccess() error {
	return d.client.Incr("athenahealth.responses.success", []string{}, 1.0)
}

func (d *Datadog) ResponseError() error {
	return d.client.Incr("athenahealth.responses.error", []string{}, 1.0)
}
