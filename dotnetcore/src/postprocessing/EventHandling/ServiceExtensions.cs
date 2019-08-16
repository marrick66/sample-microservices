using System;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using postprocessing.Events;

namespace postprocessing.EventHandling
{
    public static class ServiceExtensions
    {
        public static IServiceCollection AddEventBus(this IServiceCollection Services)
        {
            Services.AddSingleton<IEventBus, AMQPEventBus>();
            Services.AddSingleton<IEventHandler<JobRegistered>, JobRegisteredEventHandler>();
            return Services;
        }

        public static void UseEventBus(this IApplicationBuilder Builder, out IEventBus Bus)
        {
            Bus = Builder.ApplicationServices.GetService<IEventBus>();
            var handler = Builder.ApplicationServices.GetService<IEventHandler<JobRegistered>>();

            Bus.Subscribe(handler);
        }
    }
}