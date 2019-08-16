using System.Collections.Generic;
using System.Threading.Tasks;

namespace postprocessing.Models 
{
    public interface IRepository<T>
    {
        Task<T> Get(string Id);
        Task<IEnumerable<T>> GetAll();
        Task Set(T Obj);
    }
}