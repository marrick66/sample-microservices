using Microsoft.Extensions.DependencyInjection;
using postprocessing.Events;
using postprocessing.Integration.Converters;
using System;

namespace postprocessing.Integration.Handling
{
    public static class ServiceExtensions
    {
        public static IServiceCollection AddEventHandlers(this IServiceCollection Services)
        {
            Services.AddSingleton<IByteConverter, JsonByteConverter>();
            Services.AddSingleton<IObserver<JobRegistered>, JobRegisteredEventHandler>();
            return Services;
        }

    }
}
