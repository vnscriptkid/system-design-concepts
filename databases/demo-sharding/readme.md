## Setup
- Considerations?
    - Do we need to shard all tables?
    - Is that possible to have a table sharded while another not?
    - When to consider sharding?
        - CAPACITY: When the database is too large to fit in a single machine
        - LOAD (WRITEs): When the database is too slow to handle the load

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SELECT * from master_add_node('citus_worker_1', 5432);
SELECT * from master_add_node('citus_worker_2', 5432);

CREATE TABLE users (
    user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text
);

SELECT create_distributed_table('users', 'user_id');

CREATE TABLE posts (
    post_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), 
    user_id uuid REFERENCES users(user_id), 
    content text
    -- Distributed relations cannot have UNIQUE, EXCLUDE, or PRIMARY KEY constraints that do not include the partition column (with an equality operator if EXCLUDE). Fix:
    -- CONSTRAINT user_id_post_id_unique UNIQUE (user_id, post_id)
);

SELECT create_distributed_table('posts', 'user_id');
-- strategy: hash partitioning
-- all posts by the same user are guaranteed to be on the same shard

CREATE TABLE interactions (
    interaction_id uuid DEFAULT uuid_generate_v4(), 
    post_id uuid,
    user_id uuid, 
    interaction_type text CHECK (interaction_type IN ('like', 'comment', 'share')), -- Type of interaction
    content text,-- This is relevant for comments. For likes or shares, this field can be left as NULL.
    PRIMARY KEY (interaction_id, post_id, user_id),
    FOREIGN KEY (post_id, user_id) REFERENCES posts(post_id, user_id)
);

SELECT create_distributed_table('interactions', 'post_id');

```

- Considerations
    - Choosing shard key: common access patterns
    - Query by other fields than shard key: cross-shard query
        - Solution 1: secondary index
            - Pros: Faster
            - Cons: More storage, still need to query all shards
        - Solution 2: denormalization
            - Pros: Fastest, no need to query all shards
            - Cons: More storage, more complex
