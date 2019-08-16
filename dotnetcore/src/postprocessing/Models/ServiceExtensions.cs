using Microsoft.Extensions.DependencyInjection;

namespace postprocessing.Models
{
    public static class ServiceExtensions
    {
        public static IServiceCollection AddRepositories(this IServiceCollection Services)
        {
            Services.AddSingleton<IRepository<Job>, JobRepository>();
            return Services;
        }
    }
}