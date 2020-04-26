using System;
using System.Threading;
using System.Threading.Tasks;
using System.IO;
using Grpc.Core;
using Sample;
using System.Net.Http;
using Grpc.Net.Client;

namespace HelloWorldClient
{
    class Program
    {
        static async Task Main(string[] args)
        {
            var channel = new Channel("127.0.0.1:12341", ChannelCredentials.Insecure);
            var client = new HelloWorldService.HelloWorldServiceClient(channel);

            try {
                var response = await client.TestAsync(new TestRequest {Query = "random"});
                Console.WriteLine(response.Message);
            } 
            catch(RpcException e)
            {
                Console.WriteLine($"gRPC error: {e.Status.Detail}");
                Console.WriteLine($"{e}");
            }
            catch 
            {
                Console.WriteLine($"Unexpected error calling HelloWorldService");
                throw;
            }

            await channel.ShutdownAsync();
            Console.WriteLine("Press any key to exit...");
            Console.ReadKey();
        }
    }
}