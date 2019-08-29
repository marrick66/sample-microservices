using Newtonsoft.Json;
using postprocessing.Models;

namespace postprocessing.Events
{
    public class JobStatusUpdate : IEvent
    {
        [JsonIgnore]
        public string Key => "jobevents.statusupdate";
        public string ID { get; set; }
        public JobStatus Status { get; set; }
    }

}