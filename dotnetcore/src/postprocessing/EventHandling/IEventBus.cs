using System.Threading.Tasks;

namespace postprocessing.EventHandling
{
    public interface IEventBus
    {
        Task Publish(IEvent Event);
        Task Subscribe<T>(IEventHandler<T> Handler) where T: IEvent;
    }
}