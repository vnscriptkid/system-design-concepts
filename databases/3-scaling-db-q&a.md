
1. What is vertical scaling?
A. Adding more databases
B. Adding more resources like CPU, RAM, and disk to a single database
C. Splitting the database into multiple tables
D. Moving the database to the cloud
Answer: B. Adding more resources like CPU, RAM, and disk to a single database
Explanation: Vertical scaling refers to scaling up the resources of a single database server.

2. Which of the following is a limitation of vertical scaling?
A. Unlimited scaling potential
B. Physical hardware limitations
C. No downtime required
D. Decreased performance
Answer: B. Physical hardware limitations
Explanation: Vertical scaling is limited by the maximum capacity of the hardware on which the database is running.

3. What is horizontal scaling?
A. Adding more indexes to the database
B. Upgrading the database software
C. Adding more database servers
D. Normalizing the database schema
Answer: C. Adding more database servers
Explanation: Horizontal scaling, or scaling out, involves adding more nodes to the system to distribute the load.

4. What are read replicas used for in horizontal scaling?
A. To increase write performance
B. To distribute read load across multiple servers
C. To replace failed database servers
D. To backup data
Answer: B. To distribute read load across multiple servers
Explanation: Read replicas are used to handle read requests, which allows the primary server to focus on write operations.

5. In a master-replica setup, what does the master handle?
A. Only read operations
B. Both read and write operations
C. Only write operations
D. Neither read nor write operations
Answer: B. Both read and write operations
Explanation: In a master-replica setup, the master handles all write operations and can also serve read operations.

6. What is synchronous replication?
A. Data is replicated to the replica after a delay
B. Data is replicated to the replica in real-time
C. Replication occurs without any consistency
D. Replication that requires manual intervention
Answer: B. Data is replicated to the replica in real-time
Explanation: Synchronous replication ensures that data is replicated to the replica as it is written to the master, ensuring strong consistency.

7. Which of the following is a feature of asynchronous replication?
A. Zero replication lag
B. Slower writes due to waiting for replication
C. Eventual consistency
D. Immediate consistency
Answer: C. Eventual consistency
Explanation: Asynchronous replication does not replicate data in real-time, leading to eventual consistency where the replica will eventually catch up to the master.

8. What is sharding in the context of databases?
A. Encrypting database data
B. Duplicating the entire database
C. Dividing a database into smaller, more manageable pieces
D. Compressing the database to save space
Answer: C. Dividing a database into smaller, more manageable pieces
Explanation: Sharding involves splitting a database into subsets called shards, each of which can be hosted on separate servers.

9. Why is sharding used?
A. To simplify the database schema
B. To improve transaction processing speed
C. To handle larger datasets and workloads more effectively
D. To reduce the cost of the database
Answer: C. To handle larger datasets and workloads more effectively
Explanation: Sharding is used to distribute the data across multiple servers, thus handling larger datasets and improving performance.

10. What is a shard key?
A. A primary key in a shard
B. A unique identifier for a database shard
C. A key used to encrypt shard data
D. A key that determines how data is distributed across shards
Answer: D. A key that determines how data is distributed across shards
Explanation: A shard key is a field in the database that determines the distribution of data across different shards.

11. In the context of sharding, what does it mean for shards to be independent?
A. Shards do not share any data with each other
B. Shards can only contain unique data
C. Shards are not replicated
D. Shards do not communicate with each other
Answer:A. Shards do not share any data with each other
Explanation: Each shard operates independently and contains a portion of the entire dataset without overlapping daa.

12. Which of the following best describes a proxy in a sharded database system?
A. A backup server that takes over if the main shard fails
B. A server that handles all write operations
C. A component that routes requests to the appropriate shard
D. A tool for monitoring the performance of shards
Answer: C. A component that routes requests to the appropriate shard
Explanation: A proxy in a sharded system directs the incoming requests to the correct shard based on the shard key.

13. What is the primary purpose of database scaling?
A. To reduce costs
B. To increase data redundancy
C. To improve performance and handle larger loads
D. To automate database management
Answer: C. To improve performance and handle larger loads
Explanation: Database scaling is primarily done to enhance the database's ability to handle more data and provide better performance under increased loads.

14. Which of the following is NOT an advantage of asynchronous replication?
A. Faster writes
B. Strong consistency
C. Less impact on the master server
D. Some replication lag
Answer: B. Strong consistency
Explanation: Asynchronous replication does not provide strong consistency because the replica may not be immediately updated as changes occur on the master server.

15. In which scenario would vertical scaling be preferred over horizontal scaling?
A. When the database load is extremely high
B. When there are minimal increases in load
C. When the system requires high availability and fault tolerance
D. When the database is already distributed geographically
Answer: B. When there are minimal increases in load
Explanation: Vertical scaling can be simpler and more cost-effective for smaller increases in load, where the limits of hardware have not yet been reached.

16. What is the main drawback of synchronous replication?
A. Data is not replicated in real-time
B. It can lead to slower write operations
C. It does not provide any data redundancy
D. It is less secure than asynchronous replication
Answer: B. It can lead to slower write operations
Explanation: Synchronous replication can slow down writes because the system must wait for the data to be confirmed on the replica before proceeding.

17. Which of the following is crucial for API servers in a system with read replicas?
A. Knowing which database to connect to for reads
B. Performing write operations
C. Managing user authentication
D. Encrypting data before sending to the replica
Answer:A. Knowing which database to connect to for reads
Explanation: API servers must be aware of the read replicas to properly distribute read requests and not overload the master server.

18. Why might a system with sharding use a proxy?
A. To manage the replication process
B. To automatically update the shard keys
C. To handle routing requests to the appropriate shard
D. To compress data before it is stored in a shard
Answer: C. To handle routing requests to the appropriate shard
Explanation: A proxy is used to direct data requests to the correct shard, ensuring efficient data retrieval.

19. What can be a reason to configure one MySQL server as a replica of another?
A. To conduct performance testing
B. To implement sharding
C. For load balancing and fault tolerance
D. To store different types of data
Answer: C. For load balancing and fault tolerance
Explanation: Configuring a replica can help distribute the load and provide a fallback in case the master server fails.

20. What is the benefit of implementing sharding by spinning up two databases, one handling keys A-M and the other N-Z?
A. Increased security
B. Improved load distribution
C. Enhanced data compression
D. Simplified database schema
Answer: B. Improved load distribution
Explanation: Splitting the data into two shards based on keys can balance the load across servers and optimize the performance of each shard.