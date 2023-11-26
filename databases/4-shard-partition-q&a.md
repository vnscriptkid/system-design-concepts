### Review Questions on Sharding and Partitioning

1. **What is sharding in the context of databases?**
   - A) Splitting data within a single machine
   - B) Distributing data across multiple machines
   - C) Copying data to a backup server
   - D) Compressing data to save space
   - **Answer: B**
     - **Explanation:** Sharding is the process of distributing data across multiple machines, which is a common approach to scale databases horizontally.

2. **What is partitioning of data?**
   - A) Encrypting data for security
   - B) Splitting a subset of data within the same instance
   - C) Duplicating data for better performance
   - D) Organizing data into tables
   - **Answer: B**
     - **Explanation:** Partitioning refers to dividing a subset of data within the same database instance to improve manageability and performance.

3. **Which of the following is an advantage of sharding?**
   - A) Simplifies operations
   - B) Increases overall storage capacity
   - C) Reduces the need for read replicas
   - D) Makes cross-shard queries cheaper
   - **Answer: B**
     - **Explanation:** Sharding increases the overall storage capacity by distributing the data across multiple servers or instances.

4. **What is one disadvantage of sharding?**
   - A) Decreases data availability
   - B) Reduces data redundancy
   - C) Operationally complex
   - D) Limits the data size
   - **Answer: C**
     - **Explanation:** Sharding can add operational complexity due to the distributed nature of the architecture and the need to manage multiple shards.

5. **What are read replicas used for in a database system?**
   - A) To store backups of data
   - B) To handle write-heavy loads
   - C) To handle read-heavy loads
   - D) To partition data
   - **Answer: C**
     - **Explanation:** Read replicas are used to handle read-heavy loads by serving read requests without affecting the primary server handling the write operations.

6. **Which of the following is true about vertical scaling?**
   - A) It involves adding more servers
   - B) It has no limits
   - C) It refers to adding more resources to a single server
   - D) It is preferred over horizontal scaling
   - **Answer: C**
     - **Explanation:** Vertical scaling, also known as scaling up, involves adding more resources like CPU, RAM, and disk space to a single server.

7. **Horizontal partitioning of a database refers to which of the following?**
   - A) Dividing data based on rows
   - B) Dividing data based on columns
   - C) Dividing data based on tables
   - D) Dividing data based on user requests
   - **Answer: A**
     - **Explanation:** Horizontal partitioning involves dividing a table into rows, which can then be distributed across different servers or shards.

8. **Why is cross-shard querying considered expensive?**
   - A) It requires more memory
   - B) It involves complex join operations across shards
   - C) It is charged by the database provider
   - D) It needs special hardware
   - **Answer: B**
     - **Explanation:** Cross-shard querying is expensive because it often requires complex join operations across different servers, which can be resource-intensive and slow.

9. **Which strategy involves adding more instances to handle increased loads?**
   - A) Vertical Partitioning
   - B) Horizontal Scaling
   - C) Vertical Scaling
   - D) Data Compression
   - **Answer: B**
     - **Explanation:** Horizontal scaling, or scaling out, involves adding more instances or machines to handle increased loads.

10. **Vertical partitioning of a database refers to which of the following?**
    - A) Dividing data into tables
    - B) Dividing data into rows
    - C) Dividing data into columns
    - D) Dividing data into different databases
    - **Answer: C**
      - **Explanation:** Vertical partitioning involves dividing a table into columns, which can help to isolate and improve the performance of database operations.

11. **When a database is sharded, what happens to the data?**
    - A) It is duplicated across multiple servers
    - B) It is split across multiple servers
    - C) It is compressed and stored in one server
    - D) It is encrypted and stored in the cloud
    - **Answer: B**
      - **Explanation:** When a database is sharded, the data is split across multiple servers, each shard holding a portion of the total data.

12. **Which of the following best describes a shard in the context of databases?**
    - A) A backup copy of a database
    - B) A database server
    - C) A partition of a database
    - D) An encrypted segment of data
    - **Answer: C**
      - **Explanation:** A shard refers to a partition of a database where a subset of the total dataset is stored.

13. **In the slides, what is indicated by 'WPS'?**
    - A) Writes per second
    - B) Watts per second
    - C) Writes per shard
    - D) Windows per second
    - **Answer: A**
      - **Explanation:** 'WPS' stands for writes per second, which is a measure of the number of write operations a database can handle in one second.

14. **What does it mean if a database server is handling '1000 WPS'?**
    - A) It is processing 1000 write operations per second
    - B) It is processing 1000 watts per second
    - C) It is handling 1000 windows per second
    - D) It is handling 1000 web pages per second
    - **Answer: A**
      - **Explanation:** If a database server is handling '1000 WPS', it means it is capable of processing 1000 write operations per second.

15. **What is the primary reason to use sharding in a database?**
    - A) To reduce the cost of storage
    - B) To improve the security of data
    - C) To handle larger loads and increase performance
    - D) To simplify database management
    - **Answer: C**
      - **Explanation:** The primary reason for using sharding is to handle larger loads by distributing the data across multiple servers, thereby increasing the performance of the database.

16. **Which of the following would be considered a scale-up approach?**
    - A) Adding more database shards
    - B) Increasing the CPU and RAM of an existing server
    - C) Creating read replicas
    - D) Implementing better compression algorithms
    - **Answer: B**
      - **Explanation:** Scaling up, or vertical scaling, involves increasing the resources such as CPU and RAM of an existing server to handle more load.

17. **In a sharded database system, what is a potential challenge of adding more shards?**
    - A) Decreasing performance
    - B) Simplifying operations
    - C) Increasing redundancy
    - D) Managing data consistency
    - **Answer: D**
      - **Explanation:** As more shards are added, managing data consistency across them can become a challenge, requiring careful coordination and synchronization.

18. **Which of the following scenarios would benefit most from database partitioning?**
    - A) When data needs to be encrypted
    - B) When the load on the database is read-heavy
    - C) When the database size is small and manageable
    - D) When there is a need to isolate and improve performance
    - **Answer: D**
      - **Explanation:** Partitioning can greatly benefit scenarios where there is a need to isolate and improve performance by dividing the data within the same database instance.

19. **What kind of partitioning would be most appropriate for queries that access many columns but few rows?**
    - A) Horizontal partitioning
    - B) Vertical partitioning
    - C) Hybrid partitioning
    - D) No partitioning
    - **Answer: B**
      - **Explanation:** Vertical partitioning would be most appropriate for queries that access many columns but few rows, as it divides the table into columns, making such queries more efficient.

20. **If a database needs to handle 1500 WPS and a single server can only handle 1000 WPS, what strategy should be employed?**
    - A) Vertical scaling
    - B) Horizontal scaling
    - C) Data normalization
    - D) Read replica addition
    - **Answer: B**
      - **Explanation:** Horizontal scaling should be employed, which means adding more servers to handle the load, since a single server is unable to manage the increased number of write operations per second.
