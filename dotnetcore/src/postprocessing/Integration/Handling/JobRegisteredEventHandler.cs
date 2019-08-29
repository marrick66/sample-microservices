using Microsoft.Extensions.Logging;
using postprocessing.Events;
using postprocessing.Integration.Publishing;
using postprocessing.Models;
using System;
using System.Threading.Tasks;

namespace postprocessing.Integration.Handling
{
    //Once a JobRegisteredEvent is received, it should be stored in the local 
    //repository for jobs. If this were functionally complete, we'd also do some processing of the job
    //and send status events back on the event bus.
    public class JobRegisteredEventHandler : IObserver<JobRegistered>
    {
        private IEventPublisher _publisher;
        private IRepository<Job> _repository;
        private ILogger _logger;

        public JobRegisteredEventHandler(
            IEventPublisher Publisher,
            IRepository<Job> Repository, 
            ILogger<JobRegisteredEventHandler> Logger)
        {
            _publisher = Publisher ??
                throw new ArgumentNullException(nameof(Publisher));
            _repository = Repository ??
                throw new ArgumentNullException(nameof(Repository));
            _logger = Logger ??
                throw new ArgumentNullException(nameof(Logger));
        }

        public void OnCompleted()
        {
            _repository.Dispose();
        }

        public void OnError(Exception error)
        {
            _logger.LogError(error, "Error receiving event.");
        }

        public void OnNext(JobRegistered value)
        {
            if(value != null)
                _repository.SetAsync(new Job { Id = value.Id, Name = value.Name, Status = JobStatus.Registered })
                    .ContinueWith(async (t) => await MockProcessingAndEventSending(value, 10), TaskContinuationOptions.NotOnFaulted);
        }

        //As a mock of some asynchronous job processing, we just delay for some seconds,
        //update the local repository, and send a JobStatus event back to the bus.
        private async Task MockProcessingAndEventSending(JobRegistered Event, int DelaySeconds)
        {
            await Task.Delay(TimeSpan.FromSeconds(DelaySeconds));
            try
            {
                await _repository.SetAsync(new Job { Id = Event.Id, Status = JobStatus.Completed });
                await _publisher.PublishAsync(new JobStatusUpdate { ID = Event.Id, Status = JobStatus.Completed });
            }
            catch(Exception ex)
            {
                _logger.LogError(ex, "Unable to send status update for {Event}.", Event);
            }
        }
    }
}