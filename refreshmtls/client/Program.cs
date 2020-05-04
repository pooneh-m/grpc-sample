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
            string serverCa = File.ReadAllText("../creds/service.pem");
            string clientKey = File.ReadAllText("../creds/client.key");
            string clientCert = File.ReadAllText("../creds/client.pem");

            var creds = new SslCredentials(serverCa, new KeyCertificatePair(clientCert, clientKey));
            var channel = new Channel("localhost:12344", creds);
            var client = new HelloWorldService.HelloWorldServiceClient(channel);

           try {
                var response = await client.TestAsync(new TestRequest {Query = "mTLS random"});
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

            Console.WriteLine("Press any key to exit...");
            Console.ReadKey();
        }
    }
}