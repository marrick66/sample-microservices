using Dapper;
using System;
using System.Data.SqlClient;
using System.Threading.Tasks;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;

namespace postprocessing.Models 
{
    //Simple repository implementation to handle Job entities. For expediency, there's raw SQL in
    //here, but in production this would have either stored procedures or a full EntityFramework model
    //done.
    public class JobRepository : IRepository<Job>
    {
        private string _connectionString;

        public JobRepository(IConfiguration Configuration)
        {
            if(Configuration == null || Configuration.GetConnectionString("Jobs") == null)
                throw new ArgumentNullException("No Jobs connection string found.");

            _connectionString = Configuration.GetConnectionString("Jobs");
        }
        public async Task<Job> GetAsync(string Id)
        {
            using(var conn = new SqlConnection(_connectionString))
            {
                return await conn.QuerySingleAsync<Job>(
                    "SELECT Id, Name, Status FROM JobStatus WHERE Id = @Id", 
                    new { Id });
            }
        }

        public async Task SetAsync(Job Obj)
        {
            using(var conn = new SqlConnection(_connectionString))
            {
                await conn.ExecuteAsync(
                    @"IF EXISTS(SELECT Id FROM JobStatus WHERE Id = @Id) 
                        UPDATE JobStatus SET Status = @Status WHERE @Id = @Id
                    ELSE
                        INSERT INTO JobStatus(Id, Name, Status) VALUES(@Id, @Name, @Status)",
                        new { Obj.Id, Obj.Name, Obj.Status});
            }
        }

        public async Task<IEnumerable<Job>> GetAll()
        {
            using(var conn = new SqlConnection(_connectionString))
            {
                return await conn.QueryAsync<Job>("SELECT TOP 100 * FROM JobStatus");
            }
        }

        public void Dispose()
        {
            //Do nothing.
        }
    }
}