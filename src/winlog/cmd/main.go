package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/prom"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/service"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/util"
)

var (
	wg        sync.WaitGroup
	eventType = util.GetEnv("EVENT_TYPE", "sysmon")
)

func main() {
	prometheus.MustRegister(prom.ClickhouseEventsSuccess)
	prometheus.MustRegister(prom.ClickhouseEventsErrors)
	prometheus.MustRegister(prom.ClickhouseEventsTotal)
	prometheus.MustRegister(prom.KafkaEventsTotal)
	wg.Add(1)
	go prom.Metrics()
	s := service.Init()
	s.Run()

	wg.Wait()

}
