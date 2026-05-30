package main

import (
	"fmt"
	"nexus-snisid/pkg/kafka"
)

func main() {

	consumer := kafka.NewConsumer("kafka:9092", "events.risk", "gateway-group")

	go consumer.Consume(func(msg []byte) {
		fmt.Println("Received risk event:", string(msg))
	})

	select {}
}
