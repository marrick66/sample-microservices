package events

import (
	"fmt"
	"jobregistration-app/converters"
	"log"

	"github.com/streadway/amqp"
)

//AMQPEventBus is the RabbitMQ implementation of an event bus.
type AMQPEventBus struct {
	connectionString string
	connection       *amqp.Connection
	publishChannel   *amqp.Channel
	consumeChannel   *amqp.Channel
	converter        converters.ByteConverter
	subscriptions    map[subscriptionKey]chan bool
}

//subscriptionKey is used uniquely identify the channel of the subscription for closing
//at a later time.
type subscriptionKey struct {
	exchange string
	topic    string
}

//NewAMQPEventBus creates an implementation of an event bus that publishes outgoing messages, and allows
//subscriptions for incoming messages. The Converter implementation passed in determines the format
//that the message has across the wire.
func NewAMQPEventBus(connectionString string, converter converters.ByteConverter) (*AMQPEventBus, error) {
	bus := &AMQPEventBus{
		connectionString: connectionString,
		converter:        converter,
		subscriptions:    make(map[subscriptionKey]chan bool)}

	if err := bus.initChannels(); err != nil {
		return nil, err
	}

	return bus, nil
}

//initChannel is an internal helper method that creates the Connection
//and sets the default client channel to use for publishing and consuming messages.
func (bus *AMQPEventBus) initChannels() error {

	var err error

	if bus.connection == nil || bus.connection.IsClosed() {
		conn, err := amqp.Dial(bus.connectionString)

		if err != nil {
			return err
		}

		bus.connection = conn
	}

	if bus.publishChannel, err = bus.connection.Channel(); err != nil {
		return err
	}

	if bus.consumeChannel, err = bus.connection.Channel(); err != nil {
		return err
	}

	return nil
}

//Close sends messages to all subscriptions done channel to stop the event loops,
//closes the publish/consume channels, then closes the connection.
func (bus *AMQPEventBus) Close() {

	for _, done := range bus.subscriptions {
		done <- true
	}

	bus.publishChannel.Close()
	bus.consumeChannel.Close()
	bus.connection.Close()

	log.Printf("AMQP event bus closed")
}

//Publish sends a message to the bus for a topic. For the RabbitMQ implementation,
//the jobs exchange already exists on the development box, so there's no need to
//declare it here.
func (bus *AMQPEventBus) Publish(exchange string, topic string, message interface{}) error {

	if bus.connection == nil || bus.connection.IsClosed() {
		if err := bus.initChannels(); err != nil {
			return err
		}
	}

	var err error
	var body []byte

	if body, err = bus.converter.ToBytes(message); err != nil {
		return err
	}

	//According to the client documentation, this is
	//asynchronous, so no need to run it as a goroutine. Should
	//profile it to be sure, though.
	if err = bus.publishChannel.Publish(
		exchange,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType: bus.converter.ContentType(),
			Body:        body}); err != nil {
		
		log.Printf("Error sending %v to queue: %v", message, err)
	}

	log.Printf("Sent message %v to the queue.", message)
	return nil
}

//Subscribe assigns a topic to an ephemeral queue assigned to an event handler, then starts the event loop
//to handle incoming messages. There's a design flaw here that needs to be dealt with. If the channel/connection closes,
//all subscriptions are lost. So, there needs to be a way to kill all the listenForEvents goroutines, and resubscribe existing handlers.
func (bus *AMQPEventBus) Subscribe(exchange string, topic string, handler EventHandler) error {

	key := subscriptionKey{
		exchange: exchange,
		topic:    topic}

	//Check to make sure there isn't already a subscription for this exchange/topic combination:
	if _, ok := bus.subscriptions[key]; ok == true {
		return fmt.Errorf("subscription for %s/%s already exists", exchange, topic)
	}

	if bus.connection == nil || bus.connection.IsClosed() {
		if err := bus.initChannels(); err != nil {
			return err
		}
	}

	var err error
	var tempQueue amqp.Queue
	var deliveryChan <-chan amqp.Delivery

	//Setup temporary queue for this consumer:
	if tempQueue, err = bus.consumeChannel.QueueDeclare("", true, true, true, false, nil); err != nil {
		return err
	}

	//Bind the temporary queue to the exchange and topic:
	if err := bus.consumeChannel.QueueBind(tempQueue.Name, topic, exchange, false, nil); err != nil {
		return err
	}

	//Get a channel that delivers messages for the topic:
	if deliveryChan, err = bus.consumeChannel.Consume(tempQueue.Name, topic, true, true, true, false, nil); err != nil {
		return err
	}

	//create a channel for stopping the event loop goroutine
	done := make(chan bool)

	//Save the subscription locally and start the event handling loop:
	bus.subscriptions[key] = done
	go bus.listenForEvents(key, deliveryChan, done, handler)

	log.Printf("Subscribed to exchange: %s/topic: %s", exchange, topic)
	return nil
}

//listenForEvents is a handler specific event loop that waits for subscribed topic messages to arrive. If they
//can be successfully converted to the handlers default event, a goroutine is called to do the work.
func (bus *AMQPEventBus) listenForEvents(key subscriptionKey, deliveryChan <-chan amqp.Delivery, done <-chan bool, handler EventHandler) {
	for {
		select {
		case message := <-deliveryChan:
			log.Printf("Received message on %v", key)
			event := handler.DefaultEvent()
			//If the message can't be converted, just log it and skip:
			if err := bus.converter.FromBytes(message.Body, event); err == nil {
				go handler.Handle(event)
			} else {
				log.Printf("Unable to convert raw message bytes: %s.", message)
			}
		case <-done:
			log.Printf("Event loop complete for %v", key)
			return
		}
	}
}
