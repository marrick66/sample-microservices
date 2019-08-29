using postprocessing.Events;
using System.Threading.Tasks;

namespace postprocessing.Integration.Publishing
{
    public interface IEventPublisher
    {
        Task PublishAsync(IEvent Event);
    }
}
