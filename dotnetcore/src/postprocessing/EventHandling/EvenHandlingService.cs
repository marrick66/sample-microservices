using System;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace postprocessing.EventHandling
{
    public class EventHandlingService<T> : IHostedService where T: IEvent
    {
        private IEventBus _bus;
        public IEventHandler<T> _handler;
        private ILogger _logger;

        public EventHandlingService(IEventBus Bus, IEventHandler<T> Handler, ILogger<EventHandlingService<T>> Logger)
        {
            _bus = Bus ??
                throw new ArgumentNullException(nameof(Bus));
            _handler = Handler ??
                throw new ArgumentNullException(nameof(Handler));
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));
        }

        public Task StartAsync(CancellationToken cancellationToken)
        {
            return _bus.Subscribe(_handler);
        }

        public Task StopAsync(CancellationToken cancellationToken)
        {
            return _bus.Unubscribe(_handler);
        }
    }
}