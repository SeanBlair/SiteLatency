Distributed Website Measurement System
=
Description
-
Distributed measurement platform to evaluate website performance and to detect regional content variation.
Detailed assignment description: http://www.cs.ubc.ca/~bestchai/teaching/cs416_2016w2/assign4/index.html

**Summary:** The system uses a server process to coordinate several worker processes to do jobs received from a client process
through an RPC-based API where a job is one specific website target. To service a client's job, the server must 
coordinate with worker processes, and then report the measurement results back to the client.

Dependencies
-
Developed and tested with Go version 1.7.4 linux/amd64 on a Linux ubuntu 14.04 LTS machine

Running Instructions
-
- Open a terminal and clone the repo: git clone https://github.com/SeanBlair/SiteLatency.git
- Navigate to the folder: cd SiteLatency
- Run server: go run server.go localhost:1111 localhost:2222

server arguments: [worker-incoming ip:port] [client-incoming ip:port]

- Open a new terminal, navigate to the SiteLatency directory and run a Worker: go run worker.go localhost:1111

Worker arguments [ip:port that the server node is listening for workers on]

Note: the system is designed to have various worker nodes running on different machines, but it only supports one
worker running on the same machine. It was tested on MS Azure virtual machines set up in different geographical regions
worldwide. If have access to multiple machines to test the system, provide the server.go program with its machine's public 
IP as its first argument, and use that IP as the argument for any worker.go program on a different machine.

- Open a new terminal, navigate to the SiteLatency directory and run a Client with the following commands:

**Measure a website:**

go run client.go -m [server ip:port] [URI] [samples]

Instructs the system's workers to perform distributed measurements to a URI prefixed by "http://". Each worker performs 
samples number of measurements. Returns a data structure containing the IP of each worker and the min/median/max latency in 
milliseconds to retrieve URI by the worker. Additionally returns a map of booleans representing any regional difference 
between the contents of URI as perceived by each worker (empty if only one worker).

*Example* go run client.go -m http://www.surfline.com 11

**Measure workers:**

go run client.go -w [server ip:port] [samples]

Instructs the system to perform distributed measurements from each of the workers to the server. Each worker performs 
samples number of measurements consisting of the round-trip time (in milliseconds) to ping-pong a message via UDP. 
Returns a data structure containing the IP of each worker and the min/median/max round-trip latency between worker 
and the server.

*Example* go run client.go -w localhost:2222 5

