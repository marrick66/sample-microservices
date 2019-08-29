using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using postprocessing.Integration.Configuration;
using postprocessing.Integration.Converters;
using postprocessing.Integration.Reactive;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace postprocessing.Integration
{
    /// <summary>
    /// This currently implements asychronous handling of integration events via RabbitMQ.
    /// Outgoing events are concurrently queued and sent to the configured exchange, while
    /// a hot asynchronous observable of events is used for consumption of incoming ones.
    /// </summary>
    public class EventConsumingService : IHostedService
    {
        private QueueIntegrationConfiguration _config;
        private IOptions<EventHandlingOptions> _options;
        private IServiceProvider _provider;
        private IConnection _connection;
        private IByteConverter _converter;
        private ILogger _logger;
        private Dictionary<Type, IObservable<BasicDeliverEventArgs>> _sequences;
        private List<IDisposable> _subscriptions;

        public EventConsumingService(
            IConnection Connection,
            IServiceProvider Provider,
            QueueIntegrationConfiguration Configuration, 
            IOptions<EventHandlingOptions> Options,
            IByteConverter Converter,
            ILogger<EventConsumingService> Logger)
        {
            _connection = Connection ??
                throw new ArgumentNullException(nameof(Connection));
            _provider = Provider ??
                throw new ArgumentNullException(nameof(Provider));
            _config = Configuration ??
                throw new ArgumentNullException(nameof(Configuration));
            _options = Options ??
                throw new ArgumentNullException(nameof(Options));
            _converter = Converter ??
                throw new ArgumentNullException(nameof(Converter));
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));

            _sequences = new Dictionary<Type, IObservable<BasicDeliverEventArgs>>();
            _subscriptions = new List<IDisposable>();
        }

        /// <summary>
        /// Populate the base sequences from the configuration, and then register
        /// the handlers populated in the options.
        /// </summary>
        public Task StartAsync(CancellationToken cancellationToken)
        {
            var configKeys = _config.IncomingMap.Keys;
            var handlerKeys = _options.Value.Registrations.Keys;

            var missingConfigKeys = handlerKeys.Except(configKeys).Select(key => key.ToString());

            if (missingConfigKeys.Count() > 0)
                throw new InvalidOperationException($"No configuration(s) specified for the types: {string.Join(',', missingConfigKeys)}.");

            var completeKeys = configKeys.Intersect(handlerKeys);
            var registrations = _options.Value.Registrations;

            foreach (var key in completeKeys)
            {
                var sequence = Observable.Create(_connection, _config.IncomingMap[key]);

                _sequences[key] = sequence;
                _subscriptions.Add(registrations[key](_provider, sequence, _converter, _logger));
            }


            return Task.CompletedTask;
        }

        /// <summary>
        /// On stop, dispose of all of the event handler subscriptions.
        /// </summary>
        public Task StopAsync(CancellationToken cancellationToken)
        {
            foreach (var subscription in _subscriptions)
                subscription.Dispose();

            return Task.CompletedTask;
        }

       
    }
}
