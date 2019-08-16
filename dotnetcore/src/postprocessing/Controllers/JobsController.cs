using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
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
        private IRepository<Job> _repository;
        private ILogger _logger;
        
        public JobsController(IRepository<Job> Repository, ILogger<JobsController> Logger)
        {
            _repository = Repository ??
                throw new ArgumentNullException(nameof(Repository));
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));
        }

        [HttpGet]
        [SwaggerOperation(
        Summary = "Gets the status of all jobs.",
        Description = "Returns the status of all jobs in the database.",
        OperationId = "GetAllJobs" )]
        [SwaggerResponse(200, "The all jobs were returned.", typeof(Job[]))]
        public Task<IEnumerable<Job>> Get()
        {
            return _repository.GetAll();
        }
        
    }
}
