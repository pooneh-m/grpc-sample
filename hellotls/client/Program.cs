﻿using System;
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
            string serverCa = File.ReadAllText("../creds/service.pem");

            var creds = new SslCredentials(serverCa);
            var channel = new Channel("localhost:12342", creds);
            var client = new HelloWorldService.HelloWorldServiceClient(channel);

           try {
                var response = await client.TestAsync(new TestRequest {Query = "TLS random"});
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