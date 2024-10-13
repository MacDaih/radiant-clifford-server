package worker

import (
	"context"
	"log"
	"time"

	sdk "github.com/macdaih/porter_go_sdk"
)

type Consumer struct {
	client *sdk.PorterClient
	topics []string
}

func NewConsumer(client *sdk.PorterClient, topics []string) *Consumer {
	return &Consumer{
		client: client,
		topics: topics,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	log.Println("Starting Message Consumer")

	go func() {
		for {
			if err := c.client.Subscribe(ctx, c.topics); err != nil {
				log.Printf("failed to subscribe messages : %s", err.Error())
			}
			log.Println("done subscribing messages")
			time.Sleep(10 * time.Second)
		}
	}()

	<-ctx.Done()
	log.Println("worker terminating")
}
