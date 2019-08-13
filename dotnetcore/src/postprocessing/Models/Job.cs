namespace postprocessing.Models {

    //This is pretty much a duplicate of the go RPC service, but
    //if we adhere to the "owning your own data" principle of
    //microservices, it's not really an issue as long as integration
    //events are versioned and/or don't change.

    //The possible states of a single job. 
    public enum JobStatus {
        Registered = 1,
        Running = 2,
	    Failed = 3,
	    Completed = 4
    }

    //The POCO entity that represents a single job.
    public class Job {
        public int Id { get; set; }
        public JobStatus Status{ get; set; }
    }
}