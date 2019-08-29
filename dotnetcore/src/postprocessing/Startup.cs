using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json.Converters;
using postprocessing.EventHandling;
using postprocessing.Events;
using postprocessing.Integration.Handling;
using postprocessing.Models;
using Swashbuckle.AspNetCore.Swagger;

namespace postprocessing
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            services
            .AddLogging(options => options.AddConsole())
            .AddMvc()
            .AddJsonOptions(options =>
            {
                //By default, enums are returned as integers.  We want them to be strings...
                options.SerializerSettings.Converters.Add(new StringEnumConverter());
            })
            .SetCompatibilityVersion(CompatibilityVersion.Version_2_2);

            //Configure default Swagger generation:
            services.AddSwaggerGen(opt =>
            {
                opt.SwaggerDoc("v1", new Info { Title = "Sample Job Scheduler API", Version = "v1"});
                opt.EnableAnnotations();
            });

            //Add the local repositories used:
            services.AddRepositories();

            //Add event handlers
            services.AddEventHandlers();

            //Add the background event service and register the known handlers:
            services.AddEventIntegrationServices(
                options =>
                {
                    options.RegisterHandler<JobRegistered>();
                });
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IHostingEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
                app.UseHsts();
            }

            app.UseHttpsRedirection();
            app.UseSwagger();
            app.UseSwaggerUI(opt => opt.SwaggerEndpoint("/swagger/v1/swagger.json", "Job Scheduler API"));
            app.UseMvc();
        }
    }
}
