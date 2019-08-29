using System;
using System.Collections.Generic;

namespace postprocessing.Integration.Configuration
{
    /// <summary>
    /// Store the configuration for the RabbitMQ instance, with tables mapping
    /// the event type to the queue information.
    /// </summary>
    public class QueueIntegrationConfiguration
    {
        public class OutgoingConfiguration
        {
            public string EventType { get; set; }
            public string Exchange { get; set; }
            public string Routing { get; set; }
        }

        public class IncomingConfiguration
        {
            public string Exchange { get; set; }
            public string Queue { get; set; }
            public string Routing { get; set; }
            public string EventType { get; set; }
        }

        private Dictionary<Type, OutgoingConfiguration> _outgoing;
        private Dictionary<Type, IncomingConfiguration> _incoming;

        public string ConnectionString { get; set; }
        public IEnumerable<OutgoingConfiguration> OutgoingEvents { get; set; }
        public IEnumerable<IncomingConfiguration> IncomingEvents { get; set; }

        public IDictionary<Type, OutgoingConfiguration> OutgoingMap
        {
            get
            {
                if (_outgoing == null)
                    PopulateOutgoing();

                return _outgoing;
            }
        }

        public IDictionary<Type, IncomingConfiguration> IncomingMap
        {
            get
            {
                if (_incoming == null)
                    PopulateIncoming();

                return _incoming;
            }
        }

        private void PopulateIncoming()
        {
            _incoming = new Dictionary<Type, IncomingConfiguration>();

            foreach (var config in IncomingEvents)
            {
                var eventType = Type.GetType(config.EventType);

                if (_incoming.ContainsKey(eventType))
                    throw new InvalidOperationException($"Duplicate event type {config.EventType}.");

                _incoming[eventType] = config;
            }
        }

        private void PopulateOutgoing()
        {
            _outgoing = new Dictionary<Type, OutgoingConfiguration>();

            foreach (var config in OutgoingEvents)
            {
                var eventType = Type.GetType(config.EventType);

                if (_outgoing.ContainsKey(eventType))
                    throw new InvalidOperationException($"Duplicate event type {config.EventType}.");

                _outgoing[eventType] = config;
            }
        }
    }
}
