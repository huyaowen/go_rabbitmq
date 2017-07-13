package basic

import (
	"fmt"
	"go_rabbitmq"
	"go_rabbitmq/connection"
	"log"

	"github.com/streadway/amqp"
)

/*
*最简单消息发送，队列形式，无ack,无路由键，无交换机
 */

func BasicSend(sendMessageConf go_rabbitmq.MessageConfig, queueSendConf go_rabbitmq.QueueConfig) error {

	conn := connection.GetConnection()
	ch, err := conn.Channel()
	if err != nil {
		go_rabbitmq.FailOnError(err, "Failed to open a channel")
		return fmt.Errorf("Channel: %s", err)
	}
	defer ch.Close()

	//可信消息，加入确认机制
	if sendMessageConf.Confirm {
		log.Printf("enabling publishing confirms.")
		if err := ch.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

	q, err := ch.QueueDeclare(
		queueSendConf.Name,       // name
		queueSendConf.Durable,    // durable
		queueSendConf.AutoDelete, // delete when unused
		queueSendConf.Exclusive,  // exclusive
		queueSendConf.NoWait,     // no-wait
		queueSendConf.Args,       // arguments
	)
	if err != nil {
		go_rabbitmq.FailOnError(err, "Failed to declare a queue")
		return fmt.Errorf("QueueDeclare: %s", err)
	}
	res := make(chan error, 1)
	go func() {
		err = sendMsg(ch, q, sendMessageConf)
		res <- err
	}()
	return <-res
}

func sendMsg(ch *amqp.Channel, q amqp.Queue, sendMessageConf go_rabbitmq.MessageConfig) error {

	if err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		sendMessageConf.Mandatory, // mandatory
		sendMessageConf.Immediate, // immediate
		amqp.Publishing{
			ContentType:  sendMessageConf.ContentType,
			Body:         sendMessageConf.Body,
			DeliveryMode: sendMessageConf.DeliveryMode,
		},
	); err != nil {
		return fmt.Errorf("Failed to publish a message: %s", err)
	}
	return nil

}

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
