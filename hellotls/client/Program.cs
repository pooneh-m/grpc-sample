using System;
using System.Threading;
using System.Threading.Tasks;
using System.IO;
using Grpc.Core;
using Sample;
using System.Net.Http;
using Grpc.Net.Client;
using Microsoft.Extensions.Logging;

namespace HelloWorldClient
{
    class Program
    {
        static async Task Main(string[] args)
        {
            string baseAddress = "127.0.0.1";
            string serverCa = File.ReadAllText("../creds/service.pem");

            var creds = new SslCredentials(serverCa);
            var client = CreateClientWithCert("https://" + baseAddress + ":12342", creds);

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

            Console.WriteLine("Press any key to exit...");
            Console.ReadKey();
        }

        public static HelloWorldService.HelloWorldServiceClient CreateClientWithCert(
            string baseAddress,
            SslCredentials creds)
        {

            var loggerFactory = LoggerFactory.Create(logging =>
            {
                logging.AddConsole();
                logging.SetMinimumLevel(Microsoft.Extensions.Logging.LogLevel.Trace);
            });

            // Add client cert to the handler
            var handler = new HttpClientHandler();
            handler.ServerCertificateCustomValidationCallback = 
                HttpClientHandler.DangerousAcceptAnyServerCertificateValidator;

            // Create the gRPC channel
            var channel = GrpcChannel.ForAddress(baseAddress, new GrpcChannelOptions
            {
                HttpClient = new HttpClient(handler),
                LoggerFactory = loggerFactory,
                ThrowOperationCanceledOnCancellation = true,
            });

            return new HelloWorldService.HelloWorldServiceClient(channel);
        }
    }
}