using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace postprocessing.Models 
{
    public interface IRepository<T> : IDisposable
    {
        Task<T> GetAsync(string Id);
        Task<IEnumerable<T>> GetAll();
        Task SetAsync(T Obj);
    }
}