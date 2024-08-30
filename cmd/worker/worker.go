package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"webservice/internal/collector"
	tcpclient "webservice/pkg/tcp_client"

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

func (c *Consumer) Run(ctx context.Context) error {
	cerr := make(chan error, 1)

	sub := func() error {
		return c.client.Subscribe(ctx, c.topics)
	}

	go func() {
		for {
			if err := retry(10, sub); err != nil {
				cerr <- err
			}
		}
	}()

	select {
	case err := <-cerr:
		return err
	case <-ctx.Done():
		return nil
	}
}

func Process(socket, key string, collector collector.Collector, we chan error) {
	go func() {
		for {
			if time.Now().Day() >= 1 {
				if err := collector.CleanUpWithArchive(); err != nil {
					log.Printf("failed to archive records : %s", err.Error())
				}
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	we <- retry(10, func() error {
		return tcpclient.RunTCPCLient(socket, key, collector.ReadSock)
	})
}

func retry(attempts int, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		} else {
			time.Sleep(30 * time.Second)
		}
	}

	return fmt.Errorf("failed to retry after %d attempts : %w", attempts, err)
}
