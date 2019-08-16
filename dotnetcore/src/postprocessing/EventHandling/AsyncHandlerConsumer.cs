using System;
using System.Text;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using RabbitMQ.Client;

namespace postprocessing.EventHandling
{
    //This is a custom MQ consumer that receives a message with the appropriate
    //Topic and object type, then calls the associated event handler.
    public class AsyncHandlerConsumer<T> : AsyncDefaultBasicConsumer where T : IEvent
    {
        private IEventHandler<T> _handler;
        private ILogger _logger;

        public AsyncHandlerConsumer(IModel Model, IEventHandler<T> Handler, ILogger Logger)
            : base(Model)
        {
            _handler = Handler ??
                throw new ArgumentNullException(nameof(Handler));
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));
        }

        public override async Task HandleBasicDeliver(string consumerTag, ulong deliveryTag, bool redelivered, string exchange, string routingKey, IBasicProperties properties, byte[] body)
        {
            try
            {
                var bodyJson = Encoding.UTF8.GetString(body);
                var message = JsonConvert.DeserializeObject<T>(bodyJson);
                await _handler.HandleAsync(message);
            }
            catch(Exception ex)
            {
                _logger.LogError(ex, "Error handling received message.");
            }
            finally
            {
                Model.BasicAck(deliveryTag, false);
            }
        }
    }
}