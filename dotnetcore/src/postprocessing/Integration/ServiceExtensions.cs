using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using postprocessing.Integration;
using postprocessing.Integration.Configuration;
using postprocessing.Integration.Publishing;
using RabbitMQ.Client;
using System;

namespace postprocessing.EventHandling
{
    public static class ServiceExtensions
    {
        public static IServiceCollection AddEventIntegrationServices(this IServiceCollection Services, Action<EventHandlingOptions> OptionAction = null)
        {
            Services.AddSingleton(
                provider =>
                {
                    var config = provider.GetService<IConfiguration>();
                    return config.GetSection("QueueIntegrationConfiguration").Get<QueueIntegrationConfiguration>();
                });

            Services.AddSingleton(
                provider =>
                {
                    var config = provider.GetService<QueueIntegrationConfiguration>();
                    var factory = new ConnectionFactory();
                    return factory.CreateConnection(config.ConnectionString);
                });

            //Since this is used as a background service and a dependency, need to setup
            //this as a singleton multiple ways.
            Services.AddSingleton<EventPublishingService>();
            Services.AddSingleton<IEventPublisher>(provider => provider.GetService<EventPublishingService>());
            Services.AddSingleton<IHostedService>(provider => provider.GetService<EventPublishingService>());

            Services.AddSingleton<IHostedService, EventConsumingService>();

            if(OptionAction != null)
                Services.Configure(OptionAction);

            return Services;
        }
    }
}