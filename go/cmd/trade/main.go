package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/HTM1000/homebroker/go/internal/infra/kafka"
	"github.com/HTM1000/homebroker/go/internal/market/dto"
	"github.com/HTM1000/homebroker/go/internal/market/entity"
	"github.com/HTM1000/homebroker/go/internal/market/transformer"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChannel := make(chan *ckafka.Message)

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	}

	producer := kafka.NewKfakaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	go kafka.Consume(kafkaMsgChannel) // T2

	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() // T3

	go func() {
		for msg := range kafkaMsgChannel {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInputDTO := dto.TradeInputDTO{}
			err := json.Unmarshal(msg.Value, &tradeInputDTO)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInputDTO)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "   ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
