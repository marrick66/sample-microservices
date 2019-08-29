using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using postprocessing.Events;
using postprocessing.Integration.Converters;
using postprocessing.Integration.Reactive;
using RabbitMQ.Client.Events;
using System;
using System.Collections.Generic;
using System.Reactive.Concurrency;
using System.Reactive.Linq;

namespace postprocessing.Integration.Configuration
{
    /// <summary>
    /// This is a really tricky options class, not sure
    /// if I really like this implementation.
    /// </summary>
    public class EventHandlingOptions
    {
        public delegate IDisposable RegistrationFunction(
            IServiceProvider Provider, 
            IObservable<BasicDeliverEventArgs> Sequence, 
            IByteConverter Converter,
            ILogger Logger);

        public Dictionary<Type, RegistrationFunction> Registrations { get; }

        public EventHandlingOptions()
        {
            Registrations = new Dictionary<Type, RegistrationFunction>();
        }

        /// <summary>
        /// There doesn't seem to be a clean way to register different observer types without
        /// generics. I suppose we could make it the handler's responsibility, but that seems
        /// to violate SRP.  Will have to mull on this to see a better way. For now, we just create the
        /// closure and let the service call them with what it knows.
        /// </summary>
        public void RegisterHandler<T>() 
            where T : IEvent
        {
            if (Registrations.ContainsKey(typeof(T)))
                throw new InvalidOperationException($"A handler for {typeof(T)} has already been registered.");

            Registrations[typeof(T)] =
                (provider, seq, converter, logger) =>
                {
                    var handler = provider.GetService<IObserver<T>>();

                    if (handler == null)
                        throw new InvalidOperationException($"No handler for {typeof(T)} could be found.");

                    var deserializedSeq = seq.Select(args => converter.FromBytes<T>(args.Body))
                        .LogAndContinue(logger)
                        .ObserveOn(ThreadPoolScheduler.Instance);

                    return deserializedSeq.Subscribe(handler);
                    
                };
        }
    }
}
