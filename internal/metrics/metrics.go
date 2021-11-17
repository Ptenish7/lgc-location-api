package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ozonmp/lgc-location-api/internal/model"
)

var (
	locationNotFoundTotal = promauto.NewCounter(prometheus.CounterOpts{
		Subsystem: "lgc_location_api",
		Name:      "location_not_found_total",
		Help:      "Total number of locations that were not found",
	})

	eventCUDTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "lgc_location_api",
		Name:      "event_cud_total",
		Help:      "Total number of event CUD operations",
	}, []string{"event_type"})

	eventsInRetranslatorCount = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "lgc_location_api",
		Name:      "events_in_retranslator_count",
		Help:      "Number of events being processed in retranslator",
	})
)

// IncLocationNotFoundCounter increments locationNotFoundTotal counter
func IncLocationNotFoundCounter() {
	locationNotFoundTotal.Inc()
}

// IncEventCUDCounter increments eventCUDTotal counter with specified event type label
func IncEventCUDCounter(eventType model.EventType) {
	eventCUDTotal.WithLabelValues(eventType.String()).Inc()
}

// AddEventsInRetranslator adds count to eventsInRetranslatorCount gauge
func AddEventsInRetranslator(count int) {
	eventsInRetranslatorCount.Add(float64(count))
}

// SubtractEventsInRetranslator subtracts count from eventsInRetranslatorCount gauge
func SubtractEventsInRetranslator(count int) {
	eventsInRetranslatorCount.Sub(float64(count))
}
