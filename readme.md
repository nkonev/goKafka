```
go run ./cmd/producer/... --brokers=127.0.0.1:9092 --topic ololo --producers 1
go run ./cmd/consumer/... --brokers=127.0.0.1:9092 --topics ololo --group myg
```
