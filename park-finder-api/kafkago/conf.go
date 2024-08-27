package kafkago

import (
	"crypto/tls"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

func newConfig() *sarama.Config {
	conf := sarama.NewConfig()
	conf.ClientID = os.Getenv("CLOUDKARAFKA__TOPIC")
	conf.Consumer.MaxWaitTime = 2 * time.Second
	conf.Consumer.Retry.Backoff = 3
	conf.Consumer.Offsets.Initial = sarama.ReceiveTime
	conf.ChannelBufferSize = 256

	conf.Producer.Retry.Max = 3
	conf.Producer.Timeout = 2 * time.Second
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Compression = sarama.CompressionLZ4
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true

	algorithm := os.Getenv("CLOUDKARAFKA__AUTH_MECHANISM")

	isHasSHA := strings.HasPrefix(algorithm, "SCRAM-SHA-")

	if isHasSHA {
		log.Println("connecting with sha mode ...")
		conf.ClientID = "sasl_scram_client"
		conf.Metadata.Full = true
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = os.Getenv("CLOUDKARAFKA_USERNAME")
		conf.Net.SASL.Password = os.Getenv("CLOUDKARAFKA_PASSWORD")
		conf.Net.SASL.Handshake = true
		if algorithm == "SCRAM-SHA-512" {
			conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
			conf.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA512)
		} else if algorithm == "SCRAM-SHA-256" {
			conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
			conf.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA256)
		} else {
			panic("invalid SHA algorithm should be sha512 or sha256")
		}

		conf.Net.TLS.Enable = true
		conf.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: false,
		}
	}
	return conf
}
