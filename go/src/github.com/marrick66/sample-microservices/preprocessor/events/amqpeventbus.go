package events

import (
	"encoding/json"
	"os"

	"github.com/streadway/amqp"
)

//AMQPEventBus is the RabbitMQ implementation of an event bus.
type AMQPEventBus struct {
	exchange         string
	connectionString string
	contentType      string
	connection       *amqp.Connection
	channel          *amqp.Channel
}

//NewAMQPEventBus creates the implementation with defaults.
func NewAMQPEventBus() (*AMQPEventBus, error) {
	return &AMQPEventBus{
		exchange:         os.Getenv("EXCHANGE"),
		connectionString: os.Getenv("EVENT_BUS"),
		contentType:      "application/json"}, nil
}

//initChannel creates a new channel to be used internally.
func (bus *AMQPEventBus) initChannel() error {

	var err error
	var channel *amqp.Channel

	if bus.connection == nil || bus.connection.IsClosed() {
		conn, err := amqp.Dial(bus.connectionString)

		if err == nil {
			bus.connection = conn
			channel, err = bus.connection.Channel()
		}

	}

	if err != nil {
		return err
	}

	bus.channel = channel
	return nil
}

//Publish sends a message to the bus for a topic.
func (bus *AMQPEventBus) Publish(topic string, message interface{}) error {

	var err error
	if bus.connection == nil || bus.connection.IsClosed() {
		err = bus.initChannel()
	}

	if err != nil {
		return err
	}

	body, err := json.Marshal(message)

	if err == nil {
		err = bus.channel.Publish(
			bus.exchange,
			topic,
			false,
			false,
			amqp.Publishing{
				ContentType: bus.contentType,
				Body:        body})
	}

	if err != nil {
		return err
	}

	return nil
}

//Subscribe assigns a topic to an event handler, and starts its event loop.
func (bus *AMQPEventBus) Subscribe(topic string, handler EventHandler) error {

	var err error

	if bus.connection == nil || bus.connection.IsClosed() {
		err = bus.initChannel()
	}

	if err != nil {
		return err
	}

	tempQueue, err := bus.channel.QueueDeclare("", true, true, true, false, nil)

	var deliveryChan <-chan amqp.Delivery
	if err == nil {
		err = bus.channel.QueueBind(tempQueue.Name, topic, bus.exchange, false, nil)
		deliveryChan, err = bus.channel.Consume(tempQueue.Name, "", true, true, true, false, nil)
	}

	if err != nil {
		return err
	}

	go bus.listenForEvents(deliveryChan, handler)
	return nil
}

//listenForEvents is a handler specific event loop that waits for subscribed topic messages to arrive.
func (bus *AMQPEventBus) listenForEvents(deliveryChan <-chan amqp.Delivery, handler EventHandler) {
	for {
		message := <-deliveryChan

		var event interface{}
		err := json.Unmarshal(message.Body, event)

		if err == nil {
			handler.Handle(event)
		}
	}
}
