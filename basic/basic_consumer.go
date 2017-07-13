package basic

import (
	"fmt"
	"go_rabbitmq"
	"go_rabbitmq/connection"
	"log"

	"github.com/streadway/amqp"
)

func BasicReceiver(queueReceiveCong go_rabbitmq.QueueConfig, consumer go_rabbitmq.ConsumerConfig) error {
	conn := connection.GetConnection()
	go func() {
		fmt.Printf("closing: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
	}()
	ch, err := conn.Channel()
	if err != nil {
		go_rabbitmq.FailOnError(err, "Failed to open a channel")
		return fmt.Errorf("Channel: %s", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueReceiveCong.Name,       // name
		queueReceiveCong.Durable,    // durable
		queueReceiveCong.AutoDelete, // delete when unused
		queueReceiveCong.Exclusive,  // exclusive
		queueReceiveCong.NoWait,     // no-wait
		queueReceiveCong.Args,       // arguments
	)
	if err != nil {
		go_rabbitmq.FailOnError(err, "Failed to declare a queue")
		return fmt.Errorf("QueueDeclare: %s", err)
	}

	deliveries, err := ch.Consume(
		q.Name,             // queue
		consumer.Consumer,  // consumer
		consumer.AutoAck,   // auto-ack
		consumer.Exclusive, // exclusive
		consumer.NoLocal,   // no-local
		consumer.NoWait,    // no-wait
		consumer.Args,      // args
	)
	if err != nil {
		go_rabbitmq.FailOnError(err, "Failed to register a consumer")
		return fmt.Errorf("Queue Consume: %s", err)
	}
	go func() {
		consumer.Handler(deliveries, consumer.Done)
	}()
	return <-consumer.Done
}

//handler写法
func Handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}
	done <- nil
}
