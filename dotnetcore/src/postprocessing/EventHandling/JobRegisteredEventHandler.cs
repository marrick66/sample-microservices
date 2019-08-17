using System;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using postprocessing.Events;
using postprocessing.Models;

namespace postprocessing.EventHandling
{
    //Once a JobRegisteredEvent is received, it should be stored in the local 
    //repository for jobs. If this were functionally complete, we'd also do some processing of the job
    //and send status events back on the event bus.
    public class JobRegisteredEventHandler : IEventHandler<JobRegistered>
    {
        private IEventBus _bus;
        private IRepository<Job> _repository;
        private ILogger _logger;

        public JobRegisteredEventHandler(IEventBus Bus, IRepository<Job> Repository, ILogger<JobRegisteredEventHandler> Logger)
        {
            _bus = Bus ??
                throw new ArgumentNullException(nameof(Bus));
            _repository = Repository ??
                throw new ArgumentNullException(nameof(Repository));
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));
        }

        public string Topic => "jobevents.jobregistered";

        public string QueueName => "jobRegisteredQueue";

        //When the registered event is received, we add it to the local repository,
        //then do some fake processing and send the status back to the event bus.
        public async Task HandleAsync(JobRegistered Event)
        {
            _logger.LogInformation("Received event {Event}.", Event);
           try
           {
               await _repository.Set(
                   new Job 
                   { 
                       Id = Event.Id, 
                       Name = Event.Name, 
                       Status = JobStatus.Registered 
                    });
           }
           catch(Exception ex)
           {
               _logger.LogError(ex, "Unable to store registered job {Event}", Event);
               return;
           }

            //Fake some kind of processing before sending back the status event.
            await MockProcessingAndEventSending(Event, 5);
        }

        //As a mock of some asynchronous job processing, we just delay for some seconds,
        //update the local repository, and send a JobStatus event back to the bus.
        private async Task MockProcessingAndEventSending(JobRegistered Event, int DelaySeconds)
        {
            await Task.Delay(TimeSpan.FromSeconds(DelaySeconds));
            try
            {
                await _bus.Publish(new JobStatusUpdate { ID = Event.Id, Status = JobStatus.Completed });
                await _repository.Set(new Job{ Id = Event.Id, Status = JobStatus.Completed});
            }
            catch(Exception ex)
            {
                _logger.LogError(ex, "Unable to send status update for {Event}.", Event);
            }
        }
    }
}