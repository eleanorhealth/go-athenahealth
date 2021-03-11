package stats

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

const datadogDefaultNamespace = "athenahealth-client"

type Datadog struct {
	statsClient *statsd.Client
}

func NewDatadog(addr, namespace string) (*Datadog, error) {
	if len(addr) == 0 {
		addr = ":8125"
	}

	if len(namespace) == 0 {
		namespace = datadogDefaultNamespace
	}

	statsClient, err := statsd.New(addr)
	if err != nil {
		return nil, err
	}

	statsClient.Namespace = fmt.Sprintf("%s.", namespace)

	return &Datadog{
		statsClient: statsClient,
	}, nil
}

func (d *Datadog) Request() error {
	err := d.statsClient.Count("requests.count", 1, []string{}, 1.0)
	if err != nil {
		return err
	}

	return d.statsClient.Incr("requests", []string{}, 1.0)
}

func (d *Datadog) ResponseSuccess() error {
	return d.statsClient.Incr("responses.success", []string{}, 1.0)
}

func (d *Datadog) ResponseError() error {
	return d.statsClient.Incr("responses.error", []string{}, 1.0)
}
