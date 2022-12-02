package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

func main() {
	//demoProducer()
	demoConsumer()
}

// demoProducer
func demoProducer() {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{"10.248.129.188:9876"}),
		producer.WithRetry(2),
	)
	if err != nil {
		fmt.Printf("NewProducer error: %s", err.Error())
		os.Exit(1)
	}
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	topic := "TEST_TOPIC"

	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
		}
		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

func demoConsumer() {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"10.248.129.188:9876"}),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName("TEST_CONSUMER_GROUP"),
		consumer.WithAutoCommit(false),
	)
	if err != nil {
		panic(err)
	}

	err = c.Subscribe("TEST_TOPIC", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		for _, msg := range msgs {
			rlog.Info("Subscribe Callback", map[string]interface{}{
				"body": string(msg.Body),
			})
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}

	select {}
}
