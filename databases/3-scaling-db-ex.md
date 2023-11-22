# Database Scaling Practical Exercises

## Vertical Scaling Resource Monitoring
1. Set up a basic SQL or NoSQL database.
2. Load it with a substantial amount of data.
3. Perform a stress test using tools like Apache JMeter to simulate high traffic.
4. Monitor the performance metrics such as CPU, memory, and disk I/O.
5. Upgrade the server resources (e.g., CPU and RAM) and compare the performance metrics.

## Read Replica Implementation
1. Create a master database and populate it with sample data.
2. Set up a read replica of the master database.
3. Configure your application to direct all write operations to the master and read operations to the replica.
4. Measure the performance improvement in read operations.

## Synchronous Replication Setup
1. Implement synchronous replication between a master and a replica database.
2. Observe the behavior when you perform write operations—note the write latency.
3. Shut down the replica and observe the behavior of the master database to understand the implications of synchronous replication on availability.

## Asynchronous Replication Delay Measurement
1. Set up asynchronous replication between two databases.
2. Write a script to update the master database and immediately read from the replica.
3. Measure the replication lag and understand the consistency trade-offs.

## Sharding Key Selection Strategy
1. Design a sharded database system for an e-commerce platform.
2. Decide on a sharding key (e.g., customer ID, product ID).
3. Write pseudo-code for how the application will interact with different shards based on the sharding key.

## Sharding Simulation with Docker Containers
1. Use Docker to create multiple containers, each simulating a database shard.
2. Implement a simple proxy service to route requests to the correct container based on a sharding key.
3. Experiment with different sharding keys and observe the distribution of data.

## Geo-distributed Database Sharding
1. Imagine a global application that needs to be fast for users in different regions.
2. Create a sharded setup with each shard located in a different geographic region.
3. Write a script that directs users to the shard closest to them for lower latency.

## Failover Testing with Read Replicas
1. Set up a master-replica configuration.
2. Simulate a failure of the master database and promote the read replica to be the new master.
3. Observe the failover process and ensure the application continues to function.

## Load Balancing with Multiple Read Replicas
1. Create several read replicas for a heavily read-centric application.
2. Implement load balancing to distribute the read requests evenly across all replicas.
3. Analyze the system's response times and throughput as you increase the number of replicas.

## Replica and Sharding Combined Exercise
1. Combine the concepts of replication and sharding by setting up a sharded database where each shard has one or more replicas.
2. Configure the system to handle read and write operations, ensuring that writes go to the correct shard and that reads are load-balanced across the replicas.
3. Simulate a high-traffic scenario and evaluate how the combined approach handles scalability and availability.
