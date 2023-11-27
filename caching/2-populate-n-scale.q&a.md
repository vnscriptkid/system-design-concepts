## https://chat.openai.com/c/623d61fb-cccc-4cc2-9065-823c655303f2

### What is the primary purpose of a cache in the context of an API server and a database?

A) To increase the write time to the database
B) To decrease the load on the database by storing frequently accessed data
C) To replace the database
D) To monitor API usage
Correct Answer: B - The cache is used to decrease the load on the database by storing frequently accessed data, allowing for quicker data retrieval.
### What is lazy population in caching?

A) Updating the cache whenever the database is updated
B) Checking the cache first for data before accessing the database
C) Adding all database entries to the cache preemptively
D) Deleting outdated entries from the cache
Correct Answer: B - Lazy population involves checking the cache first and if the data doesn't exist, retrieving it from the database and then caching it.
### Which caching strategy involves updating the cache and database at the same time?

A) Lazy population
B) Eager population
C) Passive population
D) Active population
Correct Answer: B - Eager population involves updating both the cache and the database simultaneously to avoid cache misses.
### Which database is mentioned in the slides as being used for caching?

A) MySQL
B) Oracle
C) MongoDB
D) Redis
Correct Answer: D - Redis is mentioned as a caching solution in the slides.
### What is the benefit of using eager population for a live cricket score application?

A) It eliminates the need for a database.
B) It updates the database more frequently.
C) It saves the cache miss by proactively updating the cache.
D) It slows down the read operations for consistency.
Correct Answer: C - Eager population saves the cache miss by proactively updating the cache when data changes, which is beneficial for real-time applications like live cricket scores.
### What does sharding in the context of caching help scale?

A) Read operations
B) Write operations
C) Database size
D) API requests
Correct Answer: B - Sharding helps scale write operations by partitioning data across multiple nodes.
### Which scaling strategy involves making the cache node bigger to handle more data?

A) Vertical scaling
B) Horizontal scaling with replicas
C) Sharding
D) Load balancing
Correct Answer: A - Vertical scaling involves increasing the capacity of a single cache node.
### In the context of the slides, what does it mean to 'proactively push data to cache'?

A) To clear the cache periodically
B) To add data to the cache before it's requested by the user
C) To update the cache only after a user request
D) To remove data from the cache after it's been used
Correct Answer: B - Proactively pushing data to cache means adding data before it's requested, anticipating future need.
### What is the downside of not using caching with databases?

A) Increased cost
B) Slower API responses
C) Faster data retrieval
D) More efficient use of storage
Correct Answer: B - Not using caching can lead to slower API responses because each request may have to retrieve data from the database.
### What is a cache miss?

A) An error in the cache
B) A successful retrieval of data from the cache
C) A failure to find the requested data in the cache
D) Deleting data from the cache
Correct Answer: C - A cache miss occurs when the requested data is not found in the cache, necessitating a database query.

### Which of the following are advantages of caching? (Choose all that apply)

A) Reduced database load
B) Improved response times
C) Increased database queries
D) Decreased scalability
E) Reduced network traffic
Correct Answers: A, B, E - Caching reduces database load, improves response times, and can reduce network traffic by handling requests locally.
### Which of the following statements are true about lazy population? (Choose all that apply)

A) It requires manual intervention to update the cache.
B) It is the most popular caching strategy.
C) It involves setting an expiry for cache entries.
D) It adds data to the cache only when it's first requested.
Correct Answers: B, C, D - Lazy population is a common caching strategy, often involves setting an expiry time for cached data, and adds data to the cache upon the first request.
### What factors are important when scaling caches? (Choose all that apply)

A) The type of data being cached
B) The read/write ratio
C) The cost of cache nodes
D) The color scheme of the data
Correct Answers: A, B, C - When scaling caches, the type of data, the read/write ratio, and the cost of nodes are important considerations.
### Which actions can be part of eager population? (Choose all that apply)

A) Updating cache during each database write
B) Preemptively adding popular data to the cache
C) Waiting for a user request to populate the cache
D) Clearing the cache periodically
Correct Answers: A, B - Eager population involves updating the cache with each write and preemptively adding data that is anticipated to be in demand.
### Which strategies are involved in scaling reads? (Choose all that apply)

A) Sharding
B) Adding more cache nodes
C) Replication across multiple nodes
D) Increasing the size of each node
Correct Answers: B, C - To scale reads, adding more cache nodes and replicating data across multiple nodes are effective strategies.