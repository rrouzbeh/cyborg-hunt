package service

import (
	"strconv"
	"time"

	"github.com/prometheus/common/log"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/click"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/kafka"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/parser"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/util"
)

type Service struct {
	stopped chan struct{}

	Name          string
	FlushInterval int
	BufferSize    int
	MinBufferSize int
}

type Msg struct{}

func Init() *Service {
	s := &Service{}
	s.Name = util.GetEnv("KAFKA_CONSUMER_NAME", "symon-consumer")
	s.FlushInterval, _ = strconv.Atoi(util.GetEnv("FLUSH_INTERVAL", "20"))
	s.BufferSize, _ = strconv.Atoi(util.GetEnv("BUFFER_SIZE", "1000"))
	s.MinBufferSize, _ = strconv.Atoi(util.GetEnv("MIN_BUFFER_SIZE", "5"))

	return s
}

func KafkaConfig() *kafka.Kafka {

	k := kafka.NewKafka()
	k.Name = util.GetEnv("KAFKA_CONSUMER_NAME", "sysmon-consumer")
	k.Version = util.GetEnv("KAFKA_VERSION", "2.1.1")
	k.Oldest = util.ParseBool(util.GetEnv("KAFKA_OLDEST", "true"))
	k.Brokers = util.GetEnv("KAFKA_BROKERS", "127.0.0.1:9092")
	k.Group = util.GetEnv("KAFKA_GROUP", "cyborg-winlog")
	k.Topic = util.GetEnv("KAFKA_TOPICS", "winlogbeat")

	return k
}

func (service *Service) Run() {
	k := KafkaConfig()
	k.Init()

	if err := k.Start(); err != nil {
		panic(err)
	}

	sysmonTick := time.NewTicker(time.Duration(service.FlushInterval) * time.Second)
	sysmon_msgs := make([]parser.Sysmon, 0, service.BufferSize)
FOR:
	for {
		select {
		case msg, more := <-k.Msgs():
			if !more {
				break FOR
			}
			if len(msg) > 0 {
				s := parser.EventHandler(msg)
				sysmon_msgs = append(sysmon_msgs, s)
				if len(sysmon_msgs) >= service.BufferSize {
					click.SysmonCommit(sysmon_msgs)
					sysmon_msgs = sysmon_msgs[:0]
					sysmonTick = time.NewTicker(time.Duration(service.FlushInterval) * time.Second)
				}
			}

		case <-sysmonTick.C:
			log.Infof("Sysmon inserted: %d", len(sysmon_msgs))
			if len(sysmon_msgs) == 0 || len(sysmon_msgs) < service.MinBufferSize {
				continue
			}
			click.SysmonCommit(sysmon_msgs)
			sysmon_msgs = sysmon_msgs[:0]
		}

	}

	click.SysmonCommit(sysmon_msgs)
	service.stopped <- struct{}{}
}
