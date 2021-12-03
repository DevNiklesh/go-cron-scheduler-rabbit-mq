package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

func (conn *Conn) StartConsumer(exch, qName, rKey string, handler func(amqp.Delivery) bool) error {
	_, err := conn.Channel.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %v", err)
	}

	err = conn.Channel.QueueBind(qName, "#."+rKey+".#", exch, false, nil)
	if err != nil {
		return fmt.Errorf("queue bind: %v", err)
	}

	err = conn.Channel.Qos(0, 0, false)
	if err != nil {
		return fmt.Errorf("consumer set prefetch: %v", err)
	}

	msgs, err := conn.Channel.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %v", err)
	}

	go func() {
		for msg := range msgs {
			if handler(msg) {
				msg.Ack(false)
			} else {
				msg.Nack(false, true)
			}
		}
	}()

	return nil
}
