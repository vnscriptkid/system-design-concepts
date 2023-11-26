## generated
1. **What is sharding in the context of databases?**
2. **What is partitioning of data?**
3. **Which of the following is an advantage of sharding?**
4. **What is one disadvantage of sharding?**
5. **What are read replicas used for in a database system?**
6. **Which of the following is true about vertical scaling?**
7. **Horizontal partitioning of a database refers to which of the following?**
8. **Why is cross-shard querying considered expensive?**
9. **Which strategy involves adding more instances to handle increased loads?**
10. **Vertical partitioning of a database refers to which of the following?**
11. **When a database is sharded, what happens to the data?**
12. **Which of the following best describes a shard in the context of databases?**
13. **In the slides, what is indicated by 'WPS'?**
14. **What does it mean if a database server is handling '1000 WPS'?**
15. **What is the primary reason to use sharding in a database?**
16. **Which of the following would be considered a scale-up approach?**
17. **In a sharded database system, what is a potential challenge of adding more shards?**
18. **Which of the following scenarios would benefit most from database partitioning?**
19. **What kind of partitioning would be most appropriate for queries that access many columns but few rows?**
20. **If a database needs to handle 1500 WPS and a single server can only handle 1000 WPS, what strategy should be employed?**

## self-doubts
- shard vs partition
    - shard: DB instance
    - partition: virtual division of data
- how to move partition from one shard to another?
- horizontal partitioning vs vertical partitioning
    - horizontal: slice by rows
    - vertical: slice by columns
- rebalancing
    - when to rebalance?
        - when a shard is overloaded
    - how to rebalance?
        - move partition from one shard to another
            - how to move partition from one shard to another?
                - create a new shard
                - move partition from one shard to another
                - delete the old shard