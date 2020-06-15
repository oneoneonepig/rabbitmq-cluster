using System;
using System.Collections.Generic;
using RabbitMQ.Client;

namespace rabbitmq_client
{
    class Program
    {
        static void Main(string[] args)
        {
            var hosts = new List<string>();
            hosts.Add("10.20.131.53");
            hosts.Add("10.20.131.54");

            var factory = new ConnectionFactory()
            {
                Port = 5672,
                UserName = "admin",
                Password = "admin",
                AutomaticRecoveryEnabled = true,
                NetworkRecoveryInterval = new TimeSpan(3000)
            };
            var conn = factory.CreateConnection(hosts, ".NET Core Client");
            var channel = conn.CreateModel();

            string exchangeName = "ha-test-ex";
            string queueName = "ha-test-q";
            string routingKey = "";

            bool durable = true;

            channel.ExchangeDeclare(exchangeName, ExchangeType.Direct);
            channel.QueueDeclare(queueName, durable, false, false, null);
            channel.QueueBind(queueName, exchangeName, routingKey, null);

            byte[] messageBodyBytes = System.Text.Encoding.UTF8.GetBytes("Hello, world!");

            uint count = 0;
            // uint limit = 5000;

            Console.CancelKeyPress += delegate
            {
                channel.Close();
                conn.Close();
                Console.WriteLine("\nMessages sent: " + count);
                System.Environment.Exit(0);
            };

            Console.WriteLine("Start sending message...");
            for (uint i = 0; ; i++)
            {
                try
                {
                    channel.BasicPublish(exchangeName, routingKey, null, messageBodyBytes);
                    System.Threading.Thread.Sleep(10);
                    count++;
                }
                catch (RabbitMQ.Client.Exceptions.AlreadyClosedException)
                {
                    // Console.WriteLine("Exception: AlreadyClosedException\n" + ex);
                    Console.WriteLine("error: connection failure, retrying after 3 seconds...");
                    System.Threading.Thread.Sleep(3000);
                }
                catch (Exception ex)
                {
                    Console.WriteLine("Unhandled exception: " + ex);
                    throw ex;
                }
            }
        }
    }
}
