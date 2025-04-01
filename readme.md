```
go run ./cmd/producer/... --brokers=127.0.0.1:9092 --topic ololo --producers 1
go run ./cmd/consumer/... --brokers=127.0.0.1:9092 --topics ololo --group myg

docker exec -it kafka bash
/opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server kafka:29092 --describe --group myg --offsets
```

Kafka Sink Connector lags [always show](https://stackoverflow.com/questions/75794506/kafka-sink-connector-lags-always-show-1-even-after-processing-all-records) 1 even after processing all records
> Hard to tell exactly based on the provided config, but it could be due to Kafka transactions. There is a known issue that when using Kafka transactions, the consumer lag never reaches 0, due to the last message being a commit message, that will not be read by the consumer.
