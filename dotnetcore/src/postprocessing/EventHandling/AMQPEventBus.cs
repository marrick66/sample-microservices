using System;
using System.Collections.Generic;
using System.Text;
using System.Threading.Tasks;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using RabbitMQ.Client;

namespace postprocessing.EventHandling
{
    //This is a local wrapper around RabbitMQ client operations for exploratory testing. If we wanted to 
    //scale this out for a large number of consumers/producers, would probably look into Google Pub/Sub, 
    //Azure Service Bus, Amazon MQ, etc. Google Pub/Sub doesn't support AMQP, if we're tied to that for
    //legacy reasons.
    public class AMQPEventBus : IEventBus
    {
        private IModel _channel;
        private IConfigurationSection _config;
        private ILogger _logger;
        private IDictionary<string, List<IAsyncBasicConsumer>> _knownTopicConsumers;
        private string _exchange;
        
        public AMQPEventBus(IConfiguration Configuration, ILogger<AMQPEventBus> Logger)
        {
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));
            
            if(Configuration == null)
                throw new ArgumentNullException(nameof(Configuration));

            _config = Configuration.GetSection("AMQP");

            if(!_config.Exists())
                throw new InvalidOperationException("AMQP configuration was not found");

            _knownTopicConsumers = new Dictionary<string, List<IAsyncBasicConsumer>>();
        }

        //Sends a message to the shared exchange with a particular topic. If any
        //consumers are connected and bound to the exchange, they'll receive it if
        //it matches their topic.
        public async Task Publish(IEvent Event)
        {
            if(_channel == null || !_channel.IsOpen)
                await InitConnection();

            try
            {
                var messageJson = JsonConvert.SerializeObject(Event);
                var messageBytes = Encoding.UTF8.GetBytes(messageJson);
                
                _channel.BasicPublish(_exchange, Event.Key, false, null,  messageBytes);
            }
            catch(Exception ex)
            {
                _logger.LogError(ex, "Unable to publish event {Event}.", Event);
            }
        }

        //Sets up an ephemeral queue that will receive messages sent to the exchange
        //with the topic specified. Once the connection is closed, that queue is deleted.
        public async Task Subscribe<T>(IEventHandler<T> Handler) where T: IEvent
        {
            if(_channel == null || !_channel.IsOpen)
                await InitConnection();

            var topic = Handler.Topic.ToLower();
            if(!_knownTopicConsumers.ContainsKey(topic))
                _knownTopicConsumers[topic] = new List<IAsyncBasicConsumer>();

            _channel.QueueDeclare(Handler.QueueName, true, true, true);
            _channel.QueueBind(Handler.QueueName, _exchange, Handler.Topic, null);
        
            var consumer = new AsyncHandlerConsumer<T>(_channel, Handler, _logger);
            var tag =_channel.BasicConsume(consumer, Handler.QueueName, true, exclusive: true);

            _knownTopicConsumers[topic].Add(consumer);
        }

        private async Task InitConnection()
        {
            var factory = new ConnectionFactory();
            //This little flag caused some grief, it's not easily
            //found in the API documentation. I actually had to 
            //trace it back through the source. Not setting it
            //prevents asynchronous consumers from delivering messages.
            factory.DispatchConsumersAsync = true;
            var retryCount = _config.GetValue<int>("connectRetries");
            var uri = _config.GetValue<string>("uri");

            factory.Uri = new Uri(uri);

            for(var i = 0; i < retryCount; i++)
            {
                try
                {
                    var conn = factory.CreateConnection();
                    _channel = conn.CreateModel();
                    _exchange = _config.GetValue<string>("exchange");
                    _channel.ExchangeDeclare(_exchange, "topic", true, false, null);
                    return;
                }
                catch(Exception ex)
                {
                    _logger.LogWarning(ex, "Unable to connect to AMQP bus, retrying...");
                    await Task.Delay(TimeSpan.FromSeconds(30));
                }
            }
            
            throw new Exception("Failed to connect to the AMQP bus.");
        }
    }
}