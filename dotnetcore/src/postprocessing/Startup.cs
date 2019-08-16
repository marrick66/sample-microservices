using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.HttpsPolicy;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Newtonsoft.Json.Converters;
using postprocessing.EventHandling;
using postprocessing.Models;
using Swashbuckle.AspNetCore.Swagger;

namespace postprocessing
{
    public class Startup
    {
        private IEventBus _bus;

        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }
        public IEventBus EventBus => _bus;

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
            
            //Add the custom event bus configuration:
            services.AddEventBus();
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
            app.UseEventBus(out _bus);

        }
    }
}
