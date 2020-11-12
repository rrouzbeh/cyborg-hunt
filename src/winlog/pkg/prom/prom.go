package prom

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ClickhouseEventsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cyborg_clickhouse_events_total",
		Help: "The total number of events to insert into ClickHouse",
	},
		[]string{"events"})

	ClickhouseEventsSuccess = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cyborg_clickhouse_events_success",
		Help: "The number of events successfully inserted into ClickHouse",
	},
		[]string{"events"})

	ClickhouseEventsErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cyborg_clickhouse_events_errors",
		Help: "The number of events didn't inserted into ClickHouse",
	},
		[]string{"events"})

	KafkaConsumerTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cyborg_kafka_Consumer_total",
		Help: "The total number of events consumer read from kafka ",
	},
		[]string{"consumer_name", "topic"})

	KafkaConsumerErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cyborg_kafka_consumer_errors",
		Help: "The number of kafka consumer errors",
	},
		[]string{"topic", "consumer_group"})

	KafkaEventsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cyborg_kafka_events_total",
		Help: "The total number of events parser read from kafka ",
	},
		[]string{"logname"})
)

func Metrics() {

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
