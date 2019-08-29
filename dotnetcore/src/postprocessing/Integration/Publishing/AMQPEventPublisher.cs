using Microsoft.Extensions.Logging;
using postprocessing.Events;
using postprocessing.Integration.Converters;
using RabbitMQ.Client;
using System;
using System.Threading.Tasks;
using static postprocessing.Integration.Configuration.QueueIntegrationConfiguration;

namespace postprocessing.Integration.Publishing
{
    /// <summary>
    /// Implements publishing a message to a single exchange/routing key combination
    /// via RabbitMQ.
    /// </summary>
    public class AMQPEventPublisher : IEventPublisher, IDisposable
    {
        private OutgoingConfiguration _config;
        private IConnection _connection;
        private IModel _model;
        private IByteConverter _converter;
        private ILogger _logger;
        private bool _isDisposed = false;

        public AMQPEventPublisher(
            IConnection Connection, 
            OutgoingConfiguration Configuration, 
            IByteConverter Converter, 
            ILogger Logger)
        {
            _connection = Connection ??
                throw new ArgumentNullException(nameof(Connection));
            _config = Configuration ??
                throw new ArgumentNullException(nameof(Configuration));
            _converter = Converter ??
                throw new ArgumentNullException(nameof(Converter));
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));

            _model = _connection.CreateModel();

            //Ensure the exchange with the routing key is declared:
            _model.ExchangeDeclare(_config.Exchange, "topic", true, false, null);
        }

        #region Dispose Methods

        public void Dispose()
        {
            Dispose(true);
        }

        protected void Dispose(bool IsDisposing)
        {
            if (_isDisposed)
                return;

            if (IsDisposing)
            {
                _model.Close();
                _isDisposed = true;
            }
        }

        #endregion

        /// <summary>
        /// Execute an asynchronous continuation of Serialize->Send for an event.
        /// </summary>
        public Task PublishAsync(IEvent Event)
        {
            return Task.Run(
                () =>
                {
                    try
                    {
                        return _converter.ToBytes(Event);
                    }
                    catch (Exception ex)
                    {
                        _logger.LogError(ex, $"Error converting {Event} to bytes.");
                        throw;
                    }
                })
                .ContinueWith(t => _model.BasicPublish(_config.Exchange, _config.Routing, null, t.Result), TaskContinuationOptions.NotOnFaulted);
        }
    }
}
