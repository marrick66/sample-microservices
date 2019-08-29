using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System;
using System.Linq;
using System.Reactive.Linq;
using static postprocessing.Integration.Configuration.QueueIntegrationConfiguration;

namespace postprocessing.Integration.Reactive
{
    public static class Observable
    {
        /// <summary>
        /// Creates a filetered AMQP observable of just the delivery events so that they can be handled later.
        /// </summary>
        public static IObservable<BasicDeliverEventArgs> Create(IConnection Connection, IncomingConfiguration Configuration)
        {
            return new EventingBasicConsumerSequence(Connection, Configuration)
                .OfType<BasicDeliverEventArgs>();
        }
    }
}
