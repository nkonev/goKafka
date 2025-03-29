```
go run ./cmd/producer/... --brokers=127.0.0.1:9092 --topic ololo --producers 1
go run ./cmd/consumer/... --brokers=127.0.0.1:9092 --topics ololo --group myg

docker exec -it kafka bash
/opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server kafka:29092 --describe --group myg --offsets
```
