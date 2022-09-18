package services

import (
	"encoding/json"
	"os"
	"rabbitmq/models"

	"github.com/streadway/amqp"
)

var ConfigConn = models.ConnectionSettings{}

func CreateQueue(queueSettings models.QueueSettings) (interface{}, error) {

	channel, err := openChannel()

	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		queueSettings.Name,
		queueSettings.Durable,
		queueSettings.AutoDelete,
		queueSettings.Exclusive,
		queueSettings.NoWait,
		queueSettings.Arguments,
	)

	if err != nil {
		return queue, err
	}

	defer channel.Close()
	return queue, nil
}

func CreateConsumer(consumerSettings models.ConsumerSettings, multiple bool, requeue bool) (chan interface{}, error) {

	channel, err := openChannel()

	if err != nil {
		return nil, err
	}

	messages, err := channel.Consume(
		consumerSettings.QueueName,
		consumerSettings.Name,
		consumerSettings.AutoAck,
		consumerSettings.Exclusive,
		consumerSettings.NoLocal,
		consumerSettings.NoWait,
		consumerSettings.Arguments,
	)

	if err != nil {
		return nil, err
	}

	messagesDecoded := make(chan interface{})

	go func() {
		for message := range messages {

			messageDecoded, err := decodeMessage(message.Body)

			if err != nil {
				message.Nack(multiple, requeue)
			} else {
				message.Ack(multiple)
				messagesDecoded <- messageDecoded
			}
		}
	}()

	return messagesDecoded, nil
}

func CreateExchange(exchangeSettings models.ExchangeSettings) error {

	channel, err := openChannel()

	if err != nil {
		return err
	}

	err = channel.ExchangeDeclare(
		exchangeSettings.Name,
		exchangeSettings.Type,
		exchangeSettings.Durable,
		exchangeSettings.AutoDelete,
		exchangeSettings.Internal,
		exchangeSettings.NoWait,
		exchangeSettings.Arguments,
	)

	if err != nil {
		return err
	}

	defer channel.Close()
	return nil
}

func CreateBind(BindSettings models.BindingSettings) error {

	channel, err := openChannel()

	if err != nil {
		return err
	}

	err = channel.QueueBind(
		BindSettings.QueueName,
		BindSettings.RoutingKey,
		BindSettings.Exchange,
		BindSettings.NoWait,
		BindSettings.Arguments,
	)

	if err != nil {
		return err
	}

	defer channel.Close()
	return nil
}

func PublishMessage(publishSettings models.PublishSettings) error {

	channel, err := openChannel()

	if err != nil {
		return err
	}

	message, err := json.Marshal(publishSettings.Publish.Body)

	if err != nil {
		return err
	}

	err = channel.Publish(
		publishSettings.Exchange,
		publishSettings.RoutingKey,
		publishSettings.Mandatory,
		publishSettings.Immediate,
		amqp.Publishing{
			ContentType: publishSettings.Publish.ContentType,
			Body:        message,
		})

	if err != nil {
		return err
	}

	defer channel.Close()
	return nil
}

/*func ConsumeMessages(consumer <-chan amqp.Delivery, multiple bool, requeue bool) chan interface{} {

	messages := make(chan interface{})

	go func() {
		for mr := range consumer {

			messageDecoded, err := decodeMessage(mr.Body)

			if err != nil {
				mr.Nack(multiple, requeue)
			} else {
				mr.Ack(multiple)
				messages <- messageDecoded
			}
		}
	}()

	return messages
}*/

func decodeMessage(message []byte) (string, error) {

	var messageDecoded string
	err := json.Unmarshal(message, &messageDecoded)

	if err != nil {
		return "", err
	} else {
		return messageDecoded, nil
	}
}

func openChannel() (*amqp.Channel, error) {

	if ConfigConn.Connection == nil || ConfigConn.Connection.IsClosed() {

		_, err := openConnection()

		if err != nil {
			return nil, err
		}
	}

	channel, err := ConfigConn.Connection.Channel()

	if err != nil {
		return nil, err
	}

	return channel, nil
}

func openConnection() (*amqp.Connection, error) {

	ConfigConn.Host = os.Getenv("RABBIT_HOST")
	ConfigConn.Port = os.Getenv("RABBIT_PORT")
	ConfigConn.User = os.Getenv("RABBIT_USER")
	ConfigConn.Pass = os.Getenv("RABBIT_PASS")

	connection, err := amqp.Dial("amqp://" + "egadeduc" + ":" + "egadeduc" + "@" + "localhost" + ":" + "5672")

	if err != nil {
		return nil, err
	} else {
		ConfigConn.Connection = connection
		return ConfigConn.Connection, nil
	}
}
