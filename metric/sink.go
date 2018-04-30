package metric

import (
	"time"

	s "github.com/armon/go-metrics"
)

type Sink struct {
	s.MetricSink
}

func NewSink() *Sink {
	return &Sink{
		MetricSink: s.NewInmemSink(time.Second, 15*time.Second),
	}
}
