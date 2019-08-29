using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using postprocessing.Events;
using postprocessing.Integration.Configuration;
using postprocessing.Integration.Converters;
using postprocessing.Integration.Publishing;
using RabbitMQ.Client;
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;

namespace postprocessing.Integration
{
    /// <summary>
    /// This background service is an asynchronously sends events to the correct publisher, by type.
    /// </summary>
    public class EventPublishingService : BackgroundService, IEventPublisher
    {
        private QueueIntegrationConfiguration _config;
        private BlockingCollection<IEvent> _publishQueue;
        private Dictionary<Type, IEventPublisher> _publishers;
        private IByteConverter _converter;
        private IConnection _connection;
        private ILogger _logger;
        private Task _eventLoop;
        private CancellationTokenSource _eventLoopSource;

        public EventPublishingService(
            IConnection Connection,
            QueueIntegrationConfiguration Configuration, 
            IByteConverter Converter, 
            ILogger<EventPublishingService> Logger)
        {
            _connection = Connection ??
                throw new ArgumentNullException(nameof(Connection));
            _config = Configuration ??
                throw new ArgumentNullException(nameof(Configuration));
            _converter = Converter ??
                throw new ArgumentNullException(nameof(Converter));
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));

            _publishQueue = new BlockingCollection<IEvent>(
                new ConcurrentQueue<IEvent>());
            _publishers = new Dictionary<Type, IEventPublisher>();
            _eventLoopSource = new CancellationTokenSource();
        }

        //Send the event to the publish queue, where the loop
        //will receive it and determine which publisher to send it to.
        public Task PublishAsync(IEvent Event)
        {
            _publishQueue.Add(Event);

            return Task.CompletedTask;
        }

        /// <summary>
        /// On start, setup the AMQP connection and all the event type->publisher mappings. Then,
        /// start the event loop.
        /// </summary>
        public override Task StartAsync(CancellationToken cancellationToken)
        {
            ConfigureTypePublisherMap();

            return ExecuteAsync(_eventLoopSource.Token);
            
        }

        /// <summary>
        /// On stop, triggers the event loop to drain the publish queue and stop. Then,
        /// disposes of any publishers that are IDisposable.
        /// IDisposable.
        /// </summary>
        public override async Task StopAsync(CancellationToken cancellationToken)
        {
            _publishQueue.CompleteAdding();

            await Task.WhenAny(_eventLoop, Task.Delay(30));

            _eventLoopSource.Cancel();
            
            foreach (var publisher in _publishers.Values)
                (publisher as IDisposable)?.Dispose();
        }

        /// <summary>
        /// Continuously attempt to pull events from the blocking queue,
        /// and publish them.
        /// </summary>
        protected override Task ExecuteAsync(CancellationToken stoppingToken)
        {
            var outgoingConfig = _config.OutgoingMap;

            _eventLoop = Task.Run(
                async () => 
                {
                    while (!stoppingToken.IsCancellationRequested && !_publishQueue.IsCompleted)
                    {
                        IEvent @event = null;

                        try
                        {
                            @event = _publishQueue.Take(stoppingToken);
                        }
                        catch (Exception ex)
                        {
                            if (ex is OperationCanceledException || ex is ObjectDisposedException)
                                _logger.LogInformation("Timeout waiting for event to publish.");
                            else
                                _logger.LogError(ex, "Error getting next event from the queue.");

                            return;
                        }

                        if (@event != null && outgoingConfig.Keys.Contains(@event.GetType()))
                        {
                            var matchingPublisher = _publishers[@event.GetType()];

                            try
                            {
                                await matchingPublisher.PublishAsync(@event);
                            }
                            catch (Exception ex)
                            {
                                _logger.LogError(ex, $"Error sending {@event} to the queue.");
                            }
                        }
                    }
                }, stoppingToken);

            if (_eventLoop.IsCompleted)
                return _eventLoop;
            else
                return Task.CompletedTask;
        }

        private void ConfigureTypePublisherMap()
        {
            foreach (var outgoing in _config.OutgoingEvents)
            {
                var eventType = Type.GetType(outgoing.EventType);
                var publisher = new AMQPEventPublisher(_connection, outgoing, _converter, _logger);

                if (_publishers.ContainsKey(eventType))
                    throw new InvalidOperationException($"There is already a publisher mapped for {eventType}");

                _publishers[eventType] = publisher;
            }
        }
    }
}
