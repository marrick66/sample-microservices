using Microsoft.Extensions.Logging;
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

        /// <summary>
        /// Helper method to catch an error in the sequence, log it, and continue with the original. 
        /// It's defined recursively, because just returning the original sequence will fail on
        /// the next error received.
        /// </summary>
        public static IObservable<T> LogAndContinue<T>(this IObservable<T> Sequence, ILogger Logger)
        {
            return Sequence.Catch<T, Exception>(
                ex =>
                {
                    Logger.LogError(ex, "Error in sequence.");
                    return Sequence.LogAndContinue(Logger);
                });
        }
    }
}
