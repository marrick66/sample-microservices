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

//Publish sends a message to the bus for a topic. For the RabbitMQ implementation,
//the jobs exchange already exists on the development box, so there's no need to 
//declare it here.
func (bus *AMQPEventBus) Publish(topic string, message interface{}) error {

	var err error
	if bus.connection == nil || bus.connection.IsClosed() {
		err = bus.initChannel()
	}

	if err != nil {
		return err
	}

	body, err := json.Marshal(message)

	//According to the client documentation, this is
	//asynchronous, so no need to run it as a goroutine. Should
	//profile it to be sure, though.
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

//Subscribe assigns a topic to an ephemeral queue assigned to an event handler, then starts the event loop 
//to handle incoming messages. There's a design flaw here that needs to be dealt with. If the channel/connection closes, 
//all subscriptions are lost. So, there needs to be a way to kill all the listenForEvents goroutines, and resubscribe existing handlers.
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

	go listenForEvents(deliveryChan, handler)
	return nil
}

//listenForEvents is a handler specific event loop that waits for subscribed topic messages to arrive.
func listenForEvents(deliveryChan <-chan amqp.Delivery, handler EventHandler) {
	for {
		message := <-deliveryChan

		//Unmarshalling needs a typed implementation to be able to 
		//deserialize to, which is handler specific.
		event := handler.DefaultEvent()
		err := json.Unmarshal(message.Body, event)

		if err == nil {
			handler.Handle(event)
		}
	}
}
