### Exercise 1: Basic Redis Setup and Data Manipulation
Objective: Set up Redis on your local machine and perform basic data caching operations.

Steps:

Install Redis on your local system by following the official guide at https://redis.io/download.
Start the Redis server and ensure it's running by connecting via the Redis CLI with the command redis-cli ping. You should get a response PONG.
Use the SET command to cache a piece of data, e.g., SET user:1001 '{"name":"John", "age":30}'.
Retrieve the cached data using the GET command, e.g., GET user:1001.
Observe how the data retrieval is instantaneous, demonstrating the speed of caching.
Concepts Demonstrated: This exercise demonstrates the basic setup of a caching system and the speed of data retrieval with caching as opposed to a database query.

### Exercise 2: Cache Invalidation
Objective: Learn how to invalidate cache entries when data changes.

Steps:

Modify the data of user:1001 in the source database.
Use the DEL command in Redis to delete the outdated cache, e.g., DEL user:1001.
Recache the updated data from the database to Redis.
Access the updated data from Redis to ensure it reflects the changes.
Concepts Demonstrated: Cache invalidation is a crucial concept in caching, ensuring that the data served from the cache is not stale.

### Exercise 3: Simulating Cache Eviction Policies
Objective: Understand how different cache eviction policies affect what data is stored.

Steps:

- Configure Redis with a max memory limit by editing the redis.conf file, setting maxmemory to a low value like 2mb.
- Set maxmemory-policy to different values (allkeys-lru, volatile-lru, allkeys-random, etc.).
Populate Redis with more data than the memory limit.
- Observe which keys are evicted based on the policy.
- Concepts Demonstrated: This exercise demonstrates cache eviction policies, which are essential when dealing with limited cache sizes.

```bash
# redis.conf
maxmemory 2mb

# maxmemory-policy default is noeviction
# set maxmemory-policy to allkeys-lru
maxmemory-policy allkeys-lru

# volatile-lru vs allkeys-lru
# volatile-lru will evict keys with an expire set, allkeys-lru will evict any key

# how to set config from redis-cli
> config set maxmemory 2mb

# how get config from redis-cli
> config get maxmemory
```

### Exercise 4: Caching Database Queries
Objective: Cache expensive database queries using Redis.

Steps:

- Identify a slow-running database query in your application.
- Execute the query and cache the result in Redis with a key that uniquely identifies the query and its parameters.
- Modify the application to check the cache before hitting the database when the query is requested again.
- Measure the time difference between fetching from the cache vs. the database.
- Concepts Demonstrated: This shows how caching can reduce the load on databases and speed up data retrieval.

### Exercise 5: Implementing Cache-Aside Pattern
Objective: Implement the cache-aside pattern in a sample application.

Steps:

- Create a simple web application that retrieves data from a database.
- Before querying the database, check if the data is available in the Redis cache.
- If the cache misses, fetch the data from the database, cache it, and then return the result.
- Simulate a high number of requests and observe the performance impact.
- Concepts Demonstrated: The cache-aside pattern is a common caching pattern where the application is responsible for loading data into the cache from the database.

### Exercise 6: Cache Warm-Up
Objective: Pre-load frequently accessed data into the cache to avoid cold starts.

Steps:

Analyze your application usage patterns to identify frequently accessed data.
Write a script that pre-loads this data into Redis before peak times.
Monitor application performance during peak times to observe the effects of cache warm-up.
Concepts Demonstrated: Cache warm-up is a proactive approach to caching, which ensures that data is already in the cache when needed, thus avoiding the initial delay of populating the cache.

### Exercise 7: Session Caching
Objective: Implement session caching in a web application.

Steps:
- Configure your web application to use Redis for storing session data.
- Create a login mechanism that stores user sessions in Redis.
- Simulate user logins and subsequent requests to demonstrate that session data is being retrieved from Redis.
- Concepts Demonstrated: Session caching is a real-world use case for caching, where user session information is stored in a cache for quick access.

Demo:
- https://github.com/expressjs/session/blob/master/index.js#L518

### Exercise 8: Content Delivery Network (CDN) Usage
Objective: Use a CDN to cache static assets of a website.

Steps:

Host static assets (images, CSS, JS) on a CDN.
Update your website to reference the CDN URLs for these assets.
Use browser dev tools to observe the cache headers and response times for these assets.
Concepts Demonstrated: CDNs are a form of caching that serve static content from geographically distributed servers to decrease load times.

### Exercise 9: API Rate Limiting with Caching
Objective: Implement rate limiting on an API endpoint using Redis.

Steps:

Create an API endpoint that serves data.
Use Redis to count the number of requests made to the endpoint by each user or IP address.
Implement logic to return a rate limit error when a user exceeds the allowed number of requests within a certain time frame.
Concepts Demonstrated: This exercise demonstrates how caching can be used for rate limiting to prevent abuse of API endpoints.

### Exercise 10: Auto-Expiration of Cache Entries
Objective: Automatically expire cache entries after a certain period.

Steps:

Use the SETEX command in Redis to cache data with an expiration time, e.g., SETEX session:1001 3600 '{"token":"abc123"}'.
After the expiration time, attempt to retrieve the data and observe that it is no longer available.
Use this mechanism to automatically clear sessions or other temporary data.
Concepts Demonstrated: Auto-expiration is a feature of many caching systems that automatically removes data after a certain period, demonstrating how to handle temporary data in a cache.

For each of these exercises, the aim is to practice and reinforce the theoretical concepts presented in the slides by applying them to scenarios that simulate real-world applications of caching. By walking through these exercises, you will gain hands-on experience with caching mechanisms and understand their impact on system performance and scalability.