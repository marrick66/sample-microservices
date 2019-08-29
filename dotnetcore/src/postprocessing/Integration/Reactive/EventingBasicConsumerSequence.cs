using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System;
using static postprocessing.Integration.Configuration.QueueIntegrationConfiguration;

namespace postprocessing.Integration.Reactive
{
    /// <summary>
    /// Takes the standard event handling consumer and converts it to an observable
    /// sequence of events that can be filtered, wrapped, etc.
    /// </summary>
    public class EventingBasicConsumerSequence : IObservable<EventArgs>
    {
        private IConnection _connection;
        private IModel _model;
        private EventingBasicConsumer _consumer;
        private IncomingConfiguration _config;
        private bool _wasStarted = false;

        public EventingBasicConsumerSequence(IConnection Connection, IncomingConfiguration Configuration)
        {
            _connection = Connection ??
                throw new ArgumentNullException(nameof(Connection));
            _config = Configuration ??
                throw new ArgumentNullException(nameof(Configuration));

            _model = Connection.CreateModel();

            InitConsumer();
        }

        /// <summary>
        /// Subscribes to the observable by creating event handlers from the
        /// consumer's events to the observer's methods. The subscription disposable
        /// contains a closure that will unregister these handlers on dispose.
        /// </summary>
        public IDisposable Subscribe(IObserver<EventArgs> Observer)
        {
            var unregisterAction = RegisterHandlers(Observer);

            //The observable is cold at the beginning. When a subscription
            //occurs, we start it if it's not already running. Some thought needs
            //to be put into this logic, since the connection, model, or consumer could fail.
            //In that case, we wouldn't want to start it again.
            lock (this)
            {
                if (!_wasStarted)
                {
                    _model.BasicConsume(_config.Queue, true, _consumer);
                    _wasStarted = true;
                }
            }

            return new ActionDisposable(unregisterAction);
        }

        /// <summary>
        /// Creates and registers event handlers mapped to the appropriate observer method, 
        /// and returns a closure that unregisters them all. Depending on the scheduler, subscriptions
        /// can be made concurrently, so synchronization is needed.
        /// </summary>
        private Action RegisterHandlers(IObserver<EventArgs> Observer)
        {
            EventHandler<BasicDeliverEventArgs> received = (object sender, BasicDeliverEventArgs args) => Observer.OnNext(args);
            EventHandler<ConsumerEventArgs> registered = (object sender, ConsumerEventArgs args) => Observer.OnNext(args);
            EventHandler<ConsumerEventArgs> unregistered = (object sender, ConsumerEventArgs args) => Observer.OnCompleted();
            EventHandler<ShutdownEventArgs> shutdown = (object sender, ShutdownEventArgs args) => Observer.OnCompleted();

            lock (this)
            {
                _consumer.Received += received;
                _consumer.Registered += registered;
                _consumer.Unregistered += unregistered;
                _consumer.Shutdown += shutdown;
            }

            //The action used to unregister the handlers for this observer, used in the disposable
            //representing the subscription.  Need to verify the closure parses "this" appropriately,
            //and points to the current instance:
            return () =>
            {
                lock (this)
                {
                    _consumer.Received -= received;
                    _consumer.Registered -= registered;
                    _consumer.Unregistered -= unregistered;
                    _consumer.Shutdown -= shutdown;
                }
            };
        }

        //Setup the observer, but don't start listening to messages yet.
        private void InitConsumer()
        {
            _consumer = new EventingBasicConsumer(_model);

            //Declare and bind the configured exchange/queue:
            _model.QueueDeclare(_config.Queue, false, false, true, null);
            _model.ExchangeDeclare(_config.Exchange, "topic", true, false, null);
            _model.QueueBind(_config.Queue, _config.Exchange, _config.Routing, null);

        }
    }
}
