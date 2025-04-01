// https://github.com/IBM/sarama/blob/main/examples/txn_producer/main.go

package main

// SIGUSR1 toggle the pause/resume consumption
import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/rcrowley/go-metrics"

	"github.com/IBM/sarama"
)

// Sarama configuration options
var (
	brokers = ""
	version = ""
	topic   = ""
	verbose = false

	recordsNumber int64 = 1

	recordsRate = metrics.GetOrRegisterMeter("records.rate", nil)
)

func init() {
	flag.StringVar(&brokers, "brokers", "", "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&version, "version", sarama.MaxVersion.String(), "Kafka cluster version")
	flag.StringVar(&topic, "topic", "", "Kafka topics where records will be copied from topics.")
	flag.Int64Var(&recordsNumber, "records-number", 10000, "Number of records sent per loop")
	flag.BoolVar(&verbose, "verbose", false, "Sarama logging")
	flag.Parse()

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(topic) == 0 {
		panic("no topic given to be consumed, please set the -topic flag")
	}
}

func main() {
	keepRunning := true
	log.Println("Starting a new Sarama producer")

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.LstdFlags))

	ctx, cancel := context.WithCancel(context.Background())

	config := sarama.NewConfig()
	config.Version = version
	config.Producer.Idempotent = true
	config.Producer.Return.Errors = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Net.MaxOpenRequests = 1
	producer, err := sarama.NewAsyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		log.Panicf("Error creating producer: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				produceTestRecords(producer)
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		<-sigterm
		log.Println("terminating: via signal")
		keepRunning = false
	}

	cancel()
	producer.Close()
}

func produceTestRecords(producer sarama.AsyncProducer) {
	// Produce some records
	var i int64
	for i = 0; i < recordsNumber; i++ {
		producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder("test")}
	}

	recordsRate.Mark(recordsNumber)
}
