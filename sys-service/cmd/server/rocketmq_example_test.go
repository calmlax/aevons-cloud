package main_test

import (
	"context"
	"fmt"
	"time"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/rocketmq"

	rmqconsumer "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type OrderCreatedEvent struct {
	OrderID  string    `json:"order_id"`
	UserID   int64     `json:"user_id"`
	Amount   float64   `json:"amount"`
	CreateAt time.Time `json:"create_at"`
}

func ExampleQueue_basic() {
	cfg := config.RocketMQConfig{
		Enabled: true,
		Mode:    "standalone",
		Address: "127.0.0.1:9876",
		Producer: config.RocketMQProducerConfig{
			Group:       "sys-service-producer",
			RetryTimes:  2,
			SendTimeout: 3,
		},
		Consumer: config.RocketMQConsumerConfig{
			Group:       "sys-service-consumer",
			Concurrency: 8,
		},
	}

	client := rocketmq.NewClient(cfg)
	if err := client.Start(); err != nil {
		panic(err)
	}
	defer client.Close()

	orderQueue := rocketmq.NewQueueWithClient[OrderCreatedEvent](client, "TopicOrder", "TagCreated", 0)
	if err := orderQueue.Subscribe(func(ctx context.Context, event OrderCreatedEvent) error {
		fmt.Printf("[Consumer] 收到订单事件: orderID=%s amount=%.2f\n", event.OrderID, event.Amount)
		return nil
	}); err != nil {
		panic(err)
	}

	event := OrderCreatedEvent{
		OrderID:  "ORD-20240409-001",
		UserID:   10086,
		Amount:   299.99,
		CreateAt: time.Now(),
	}
	result, err := orderQueue.Enqueue(context.Background(), event)
	if err != nil {
		fmt.Println("发送失败:", err)
		return
	}
	fmt.Println("[Producer] 同步发送成功, msgId:", result.MsgID)
}

func ExampleSendSync() {
	cfg := config.RocketMQConfig{
		Enabled: true,
		Mode:    "standalone",
		Address: "127.0.0.1:9876",
		Producer: config.RocketMQProducerConfig{
			Group: "sys-service-producer",
		},
	}

	client := rocketmq.NewClient(cfg)
	if err := client.Start(); err != nil {
		panic(err)
	}
	defer client.Close()

	body := []byte(`{"event":"user_registered","uid":123}`)
	result, err := client.SendSync(context.Background(), "TopicUser", "TagRegister", body)
	if err != nil {
		fmt.Println("发送失败:", err)
		return
	}
	fmt.Println("[Producer] 低层同步发送, msgId:", result.MsgID)
}

func ExampleRegisterHandler() {
	cfg := config.RocketMQConfig{
		Enabled: true,
		Mode:    "standalone",
		Address: "127.0.0.1:9876",
		Consumer: config.RocketMQConsumerConfig{
			Group:       "sys-service-consumer",
			Concurrency: 4,
		},
	}

	client := rocketmq.NewClient(cfg)
	if err := client.Start(); err != nil {
		panic(err)
	}
	defer client.Close()

	if err := client.RegisterHandler("TopicUser", "TagRegister", 4,
		func(ctx context.Context, msgs ...*primitive.MessageExt) (rmqconsumer.ConsumeResult, error) {
			for _, msg := range msgs {
				fmt.Printf("[Consumer] 低层收到消息: msgId=%s body=%s\n", msg.MsgId, string(msg.Body))
			}
			return rmqconsumer.ConsumeSuccess, nil
		},
	); err != nil {
		panic(err)
	}
}

func Example_transactionMessage() {
	cfg := config.RocketMQConfig{
		Enabled: true,
		Mode:    "cluster",
		Addresses: []string{
			"10.0.0.11:9876",
			"10.0.0.12:9876",
		},
		Producer: config.RocketMQProducerConfig{
			Group: "sys-service-tx-producer",
		},
	}

	client := rocketmq.NewClient(cfg)
	client.SetTransactionListener(rocketmq.NewFuncTransactionListener(
		func(msg *primitive.Message) primitive.LocalTransactionState {
			return primitive.CommitMessageState
		},
		func(msg *primitive.MessageExt) primitive.LocalTransactionState {
			return primitive.CommitMessageState
		},
	))
	if err := client.Start(); err != nil {
		panic(err)
	}
	defer client.Close()

	result, err := client.SendTransaction(context.Background(), "TopicOrder", "TagCreated", []byte(`{"order_id":"ORD-1"}`))
	if err != nil {
		fmt.Println("发送事务消息失败:", err)
		return
	}
	fmt.Println("[Producer] 事务消息发送成功, state:", result.State)
}
