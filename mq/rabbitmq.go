package mq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shpdwx/web/conf"
)

func NewRabbitMQ(r conf.RabbitMQ) *amqp.Channel {

	var (
		user   = r.User
		passwd = r.Passwd
		ep     = r.Endpoint
		vh     = r.Vhost
	)

	// connection rabbitmq
	dsn := fmt.Sprintf("amqp://%s:%s@%s/%s", user, passwd, ep, vh)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		fmt.Printf("Failed to connection %v", err)
		return nil
	}
	// defer conn.Close()

	// start channel
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to start a channel %v", err)
		return nil
	}
	// defer ch.Close()

	return ch
}

func LogMq(ctx context.Context, ch *amqp.Channel, r conf.RabbitMQ, s string) {

	var (
		exchange = r.Exchange
		et       = r.ExchangeType
		rt       = r.RouteKey
	)

	// declare exchange
	err := ch.ExchangeDeclare(exchange, et, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Failed to declare an exchange %v", err)
		return
	}

	err = ch.PublishWithContext(ctx, exchange, rt, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(s),
	})
	if err != nil {
		fmt.Printf("Failed to publish log %v", err)
		return
	}
}
