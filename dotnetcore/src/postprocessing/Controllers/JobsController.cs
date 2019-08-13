using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using postprocessing.Models;
using Swashbuckle.AspNetCore.Annotations;

namespace postprocessing.Controllers
{
    //This is the single default job controller that returns the status
    //of all jobs. Currently, it's open to all requests. It'll be locked down
    //in later phases, when cloud provider(s) are integrated.
    [Route("api/[controller]")]
    public class JobsController : ControllerBase
    {
        [HttpGet]
        [SwaggerOperation(
        Summary = "Gets the status of all jobs.",
        Description = "Currently returning static data on two dummy jobs.",
        OperationId = "GetAllJobs" )]
        [SwaggerResponse(200, "The all jobs were returned.", typeof(Job[]))]
        public Task<Job[]> Get()
        {
            //For the initial spike, just return a static set of data.
            return Task.FromResult(
                new Job[] 
                { 
                    new Job { Id = 1, Status = JobStatus.Registered }, 
                    new Job { Id = 2, Status = JobStatus.Completed }, 
                });
        }
        
    }
}
