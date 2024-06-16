package metric

import (
	"errors"

	"github.com/DataDog/datadog-go/statsd"
)

//NewStatsdMonitoring creates new statsd monitoring instance, datadog as default
func NewStatsdMonitoring(host string) (StatsdMonitoring, error) {
	client, err := statsd.New(host)
	if err != nil {
		return nil, err
	}
	return &datadogStatsd{client: client}, nil
}

//StatsdMonitoring contracts
type StatsdMonitoring interface {
	Increment(name string, tags []string, rate float64) error
}

type datadogStatsd struct {
	client *statsd.Client
}

func (d *datadogStatsd) Increment(name string, tags []string, rate float64) error {
	if d.client == nil {
		return errors.New("client is not initialized")
	}
	return d.client.Incr(name, tags, rate)
}
