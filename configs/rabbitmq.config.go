package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var AMQPCON *amqp.Connection
var MOBILE_DEVICE_XC = "mobile_devices"
var NOTIFICATION_XC = "notification"

func InitRabbitMQ() *amqp.Connection {
	DSN := os.Getenv("RABBITMQ_URL")
	con, err := amqp.Dial(DSN)
	if err != nil {
		log.Fatal(fmt.Sprintf("RABBITMQ: %s", err))
	}

	AMQPCON = con

	initExchanges()

	return AMQPCON
}

func initExchanges() {
	ch, err := AMQPCON.Channel()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to Init Exchanges: %s", err))
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		MOBILE_DEVICE_XC, // name
		"direct",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to Declare Exchanges: %s", err))
	}

	err = ch.ExchangeDeclare(
		NOTIFICATION_XC, // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to Declare Exchanges: %s", err))
	}
}

func InitChannel() (*amqp.Channel, error) {
	ch, err := AMQPCON.Channel()
	if err != nil {
		log.Printf("[Queue] Failed to init channel: %v", err)
		return nil, err
	}
	return ch, nil
}

func InitMobileQueue(ch *amqp.Channel, userID string, deviceID string) (*amqp.Queue, context.Context, context.CancelFunc, error) {
	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s_%s", userID, deviceID),
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// bind queue to the exchange
	if err := ch.QueueBind(q.Name, q.Name, MOBILE_DEVICE_XC, false, nil); err != nil {
		return nil, nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}
	return &q, ctx, cancel, nil
}

func InitNotifQueue(ch *amqp.Channel, userID string) (*amqp.Queue, context.Context, context.CancelFunc, error) {
	q, err := ch.QueueDeclare(
		"notification",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// bind queue to the exchange
	if err := ch.QueueBind(q.Name, q.Name, NOTIFICATION_XC, false, nil); err != nil {
		return nil, nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}
	return &q, ctx, cancel, nil
}

func SendToMQ(ch *amqp.Channel, q *amqp.Queue, ctx context.Context, payload []byte) error {
	err := ch.PublishWithContext(
		ctx,
		MOBILE_DEVICE_XC, // exchange
		q.Name,           // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func SendToNotifMQ(ch *amqp.Channel, q *amqp.Queue, ctx context.Context, payload []byte) error {
	err := ch.PublishWithContext(
		ctx,
		NOTIFICATION_XC, // exchange
		q.Name,          // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
