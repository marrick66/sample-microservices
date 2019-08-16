using System.Threading.Tasks;

namespace postprocessing.EventHandling
{
    public interface IEventHandler<T> where T: IEvent
    {
        string Topic{ get;}
        string QueueName { get; }
        Task HandleAsync(T Event);
    }

}