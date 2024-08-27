package kafkago

import "os"

const (
	DefaultTopicEnum = "park-finder-group"

	StrategyRoundRobinEnum = "roundrobin"
	StrategyStickyEnum     = "sticky"
	StrategyRangeEnum      = "range"

	GroupMainLogEnum = "park-finder-group"
)

type settingKafka struct {
	DefaultTopicEnum       string
	StrategyRoundRobinEnum string
	StrategyStickyEnum     string
	StrategyRangeEnum      string
	GroupMainLogEnum       string
}

func SettingKafkaConsumer() *settingKafka {

	prefixKafka := os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX")

	defaultTopic := os.Getenv("CLOUDKARAFKA__TOPIC")

	defaultTopic = prefixKafka + defaultTopic

	defaultGroup := os.Getenv("CLOUDKARAFKA_GROUPS")
	if defaultGroup == "" {
		defaultGroup = "azemrndm-park-finder"
	}

	defaultGroup = prefixKafka + defaultGroup

	return &settingKafka{
		DefaultTopicEnum:       defaultTopic,
		StrategyRoundRobinEnum: "roundrobin",
		StrategyStickyEnum:     "sticky",
		StrategyRangeEnum:      "range",
		GroupMainLogEnum:       defaultGroup,
	}
}
