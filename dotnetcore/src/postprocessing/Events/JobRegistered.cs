namespace postprocessing.Events
{
    public class JobRegistered : IEvent
    {
        public string Key => "jobevents.registered";
        public string Id { get; set; }
        public string Name { get; set; }
    }
}