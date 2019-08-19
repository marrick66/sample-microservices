# Sample-Microservices
This is a very rudamentary implementation of a unspecified job registration and scheduling set of microservices. Each microservice uses a different set of technologies to demonstrate the flexibility of integration. Could be relevant in the case that you want to start some custom batch learning jobs not implemented via a cloud provider's product.

## Use Case:
We assume that there are some predefined job definitions somewhere.

* __User__:
  1. Submits a request to the registration service to start a new job for the specified name.
  2. On success, receives a unique ID for the job registered.
  3. Can query the registration service to see the last known status of the job.
* __Admin__:
  1. Can query the scheduler service to see the status of all running jobs.

## Operation flow:
We assume that there are some predefined job definitions somewhere, and they are idempotent.

1. Users submit the name of the job to the registration service.
2. The registration service stores the job, and asynchronously publishes the registered event.
3. The scheduler service eventually receives the event and stores it.
4. The job is scheduled and run asynchronously at some point (mocked in the project).
5. Status updates are published to the event bus.
6. The registration service eventually receives status events and updates it's local store for the job.

## Component Services:

### Job Registration:
* gRPC service front end with operations defined in protobuf and implemented in Go.
* MongoDB for document storage of registered jobs.
* A generic event bus client implementation using RabbitMQ.
* Asynchronous event handling of messages from the bus.

### Job Scheduling/Processing:
* .Net Core ASP.Net WebAPI front end with REST resources and Swagger API documentation and implemented in C#.
* .SQL Server for storing and updating events.
* A generic event bus client implementation using RabbitMQ.
* Asychronous event handling of messages from the bus.

Both are containerized for deployment by Kubernetes or some other orchestrator service.

### An example within a cloud provider:
![Simple Diagram](./microservice.png)

## Some prerequisite(important) comments:
This is by no means production-worthy code. It represents just a few hours investment in learning and integrating Go into a microservice environment and getting familiar with VSCode. It doesn't include unit, integration, or performance tests. Assuming that eventually something like this would be implemented in production there are a plethora of things that would need to be redesigned. Since there are two separate communication methods (protobuf and REST), an API gateway should be used for external use. Versioning, request throttling, authentication, etc. would be implemented here. Assuming that's in place, it's missing any authorization mediation locally or TLS on the endpoints. No production application log aggregation or consolidating exists, and health check operations need to be implemented for load balancing or Kubernetes monitoring. From online hearsay, RabbitMQ is unreliable for cloud-scale deployments, but I don't have any empirical data to back that up. The ephemeral nature of the queues and topics means that there's only a "at most once" guarantee on events. This would need to be evaluated to see if a stronger guarantee is needed. On another note, this also represents my first exposure to the Go ecosystem, so it's probably not idiomatic code. The error handling isn't clean, and should be refactored after I've had a chance to research best practices.

### Simple scalability and fault modeling:
For simplicity, I assume that the actual execution of registered jobs is out of scope. We'll model the requirements with the following assumptions:

* All consumers of the service are located in a single geographic location, so no need for dynamic DNS to the 
closest data center.
* The ratio of reads to writes is approximately 100:1.
* Consumers expect jobs to take a non-trivial amout of time, so eventual consistency is not an issue.

#### Fault modeling:

* __Event handling__: Due to container crashes, migration, or maintenance, events will need an "at least once" guarantee on delivery. There are two possible methods that come to mind, persisting messages on the bus or having the service occasionally resend events for jobs that haven't been updated. This decision probably rests on the constraints of the service bus implementation used. All of the cloud providers have persistence as a feature, if that's where it's deployed.
* __Data consistency__: MongoDB has client defined strong or eventual consistency to replicas. Not requiring strong consistency would probably not be an issue for the use cases above, since the user can simply retry.
* __Container orchestration__: The load balancer and orchestrator will use the services health checks to determine liveness. Deployments should use rollout to maximize uptime.

#### Scalability Testing Plan (In Progress):
Since this is a personal project, there are obviously cost constraints that keep me from doing a full testing cycle with a large set of virtual users.  That being said, I am going to setup a CI/CD pipeline to what Azure resources I can get. With a personal project there's no performance criteria per se, but we should be able to get a decent estimation of what we can handle. Starting from scratch, the strategy is:

* Run profiling on small concurrent workloads (100, 500, 1000) locally in development for getting baseline of the hot paths in the code, basic CPU, memory, and I/O as well as identifying possible opportunities to optimize and/or add concurrency. 
  1. The Job Registration service will be profiled with pprof.
  2. The Job Scheduling service will be profiled with MS Perfview and MiniProfiler.

* Use the (imperfect) data along with the stated SLA response times and resource assignment for each Azure product to model usage and response times per call. This will give a t-shirt size of what we should expect to see.

*  TODO: Research and choose the Azure monitoring solution.  Application insights so far doesn't seem to support Go natively. May have to move a step ahead and deploy these as K8 pods and use that monitoring.