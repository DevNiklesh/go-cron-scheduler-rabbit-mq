package rabbit

import (
	"time"

	"github.com/streadway/amqp"
)

func (conn *Conn) Publish(exch, rKey string, message []byte) error {
	return PublishInChannel(conn.Channel, exch, rKey, message)
}

func PublishInChannel(ch *amqp.Channel, exch, rkey string, message []byte) error {
	return ch.Publish(
		exch,
		rkey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         message,
		},
	)
}
