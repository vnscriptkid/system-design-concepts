# Indexing

## How tables are stored
- Heap
- Clustered

## Indexes

## Doubts
- `CREATE INDEX` vs `CREATE UNIQUE INDEX`:
    - `CREATE INDEX`:
        - Allow duplicate values
        - Faster
    - `CREATE UNIQUE INDEX`: (primary key)
        - Don't allow duplicate values
        - Slower
- Index becomes worse when:
    - High cardinality: query returns a lot of matching rows (high percentage of total rows)


## Best practices (postgres)
#### INDEX ON `release_year`
- Applied queries:
    - `SELECT * FROM movies WHERE release_year = 2019`
    - `SELECT * FROM movies WHERE release_year = 2019 ORDER BY release_year LIMIT 100`
- Demo
```sql
-- Create table
CREATE TABLE movies (
    movie_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    release_year INTEGER NOT NULL
);

-- Seed data
-- Generate 10.000.000 records random data year from 1900 to 2020
INSERT INTO movies (title, release_year)
SELECT md5(random()::text), 1900 + (random() * 120)::integer
FROM generate_series(1, 10000000);

-- Create index
CREATE INDEX idx_movies_release_year ON movies (release_year);

-- Get info about index
SELECT * FROM pg_indexes WHERE tablename = 'movies';

-- Get size of index
SELECT pg_size_pretty(pg_indexes_size('movies'));

-- Query
EXPLAIN ANALYZE SELECT * FROM movies WHERE release_year = 2019;
-- bitmap index scan -> create a bitmap of all the rows that match the condition -> bitmap heap scan
-- experiment on high cardinality:
-- num of matching rows: 83238
-- seq scan: 1200ms
-- index scan: 1500ms

-- num of matching rows: 10
-- seq scan: 300ms
-- index scan: 0.3ms
```


#### INDEX ON `title`
- Applied queries:
    - `SELECT * FROM movies WHERE title = 'Avengers: Endgame'`
- Demo
```sql
-- Create table
CREATE TABLE movies (
    movie_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    release_year INTEGER NOT NULL
);

-- Seed data
INSERT INTO movies (title, release_year) VALUES ('Avengers: Endgame', 2019);
INSERT INTO movies (title, release_year) VALUES ('Avengers: Infinity War', 2018);
INSERT INTO movies (title, release_year) VALUES ('Black Panther', 2018);
INSERT INTO movies (title, release_year) VALUES ('Avengers: Age of Ultron', 2015);
INSERT INTO movies (title, release_year) VALUES ('The Avengers', 2012);

-- Create index
CREATE INDEX idx_movies_title ON movies (title);

-- Query
EXPLAIN ANALYZE SELECT * FROM movies WHERE title = 'Avengers: Endgame';
-- without index: seq scan -> 300ms
-- with index: index scan -> 0.1ms

```
#### INDEX ON `release_year` AND `title`
- Applied queries:
    - `SELECT * FROM movies WHERE release_year = 2019 AND title = 'Avengers: Endgame'`
- Demo
```sql
-- Seed data
INSERT INTO movies (title, release_year) VALUES ('Avengers: Endgame', 2019);
INSERT INTO movies (title, release_year) VALUES ('Avengers: Infinity War', 2018);
INSERT INTO movies (title, release_year) VALUES ('Black Panther', 2018);

-- Create index
CREATE INDEX idx_movies_release_year_title ON movies (release_year, title);

-- Query
EXPLAIN ANALYZE SELECT * FROM movies WHERE release_year = 2019 AND title = 'Avengers: Endgame';
```

#### INDEX ON  (`first_name`, `last_name`, `country`)
- Applied queries:
    - `SELECT * FROM users WHERE first_name = 'John' AND last_name = 'Doe' AND country = 'US'`
    - `SELECT * FROM users WHERE first_name = 'John' AND last_name = 'Doe'`
    - `SELECT * FROM users WHERE first_name = 'John'`
    - `SELECT * FROM users WHERE first_name = 'John' AND country = 'US'`; 
        - interesting!!! -> index can still be used (partially) even field `last_name` is skipped
        - not as efficient as INDEX ON `first_name` AND `last_name` AND `country`
        - but still better than single INDEX ON `first_name`
- Not applied queries:
    - `SELECT * FROM users WHERE last_name = 'Doe'`
    - `SELECT * FROM users WHERE country = 'US'`
    - `SELECT * FROM users WHERE last_name = 'Doe' AND country = 'US'`

```sql
-- Create table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    country TEXT NOT NULL
);

-- Seed data 10.000.000 records
INSERT INTO users (first_name, last_name, country)
SELECT md5(random()::text), md5(random()::text), md5(random()::text)
FROM generate_series(1, 10000000);
-- Insert some deterministic data
INSERT INTO users (first_name, last_name, country) VALUES 
    ('John', 'Doe', 'US'),
    ('John', 'Doe', 'UK'),
    ('John', 'Doe', 'FR'),
    ('John', 'Doe', 'IT'),
    ('Steve', 'Jobs', 'US'),
    ('Steve', 'Jobs', 'UK'),
    ('Steve', 'Jobs', 'FR'),
    ('Steve', 'Jobs', 'IT'),
    ('Steve', 'Wozniak', 'US'),
    ('Steve', 'Wozniak', 'UK'),
    ('Steve', 'Wozniak', 'FR'),
    ('Steve', 'Wozniak', 'IT');

-- Create index
CREATE INDEX idx_users_first_name_last_name_country ON users (first_name, last_name, country);

-- Create some overlapping indexes
CREATE INDEX idx_users_first_name_last_name ON users (first_name, last_name);
CREATE INDEX idx_users_first_name ON users (first_name);
```

#### Overlapping indexes
```sql
-- Create index
CREATE INDEX idx_users_first_name_last_name_country ON users (first_name, last_name, country);
CREATE INDEX idx_users_first_name_last_name ON users (first_name, last_name);
CREATE INDEX idx_users_first_name ON users (first_name);

-- Get size of each index
SELECT * FROM 
    (
        VALUES 
            ('idx_users_first_name_last_name_country', pg_size_pretty(pg_relation_size('idx_users_first_name_last_name_country'))), 
            ('idx_users_first_name_last_name', pg_size_pretty(pg_relation_size('idx_users_first_name_last_name'))), 
            ('idx_users_first_name', pg_size_pretty(pg_relation_size('idx_users_first_name')))
    ) AS t (a, b);
--                    a                    |    b    
-- ----------------------------------------+---------
--  idx_users_first_name_last_name_country | 1251 MB
--  idx_users_first_name_last_name         | 906 MB
--  idx_users_first_name                   | 563 MB
```
- Cons
    - Redundant as `idx_users_first_name_last_name` is a subset of `idx_users_first_name_last_name_country`
    - Larger index size
    - Slower write

#### Query on range
```sql
-- create tables users: id, name, active, signup_date
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL,
    signup_date DATE NOT NULL
);

-- camparision: INDEX ON (name, signup_date, active) vs INDEX ON (name, active, signup_date)
-- query: SELECT * FROM users WHERE name = 'John' AND signup_date BETWEEN '2019-01-01' AND '2019-12-31' AND active = true

INSERT INTO users (name, active, signup_date)
SELECT
    (ARRAY['John','Steve','Bill','Mark','Elon','Jeff','Larry','Sergey','Jack','Tim','Satya','Sundar','Susan','Sheryl','Marissa','Meg','Ginni','Indra','Urs','Reed','Brian','Daniel','Chuck','Bob','Tony','Travis','Dara','Ryan','Adam','Evan','Kevin','Arash','Chad','Brian','Drew','Alex','Dustin','Kevin','Jason','Aaron','Joe','Jeremy','Nathan','Julian','Dave','Ben','Matt','David','Paul','Michael','Andrew','Chris','Scott','Eric','Greg','Josh','Mike','Nick','Peter','Brian','Erik','Jed','Joshua','Max','Tom','Sam','Dan','Ian','Zach','Charlie','Rob','Ross','Gary','Chris','Steve','Dylan','Rahul','Jared','Kurt','Bobby','Trevor','Joel','Jon','Shawn','Alan','Tyler','Ethan','Brandon','Megan','Rachel','Kim','Michelle','Emily','Jessica','Emma','Lauren','Sarah','Kayla','Anna','Jennifer','Amanda','Nicole','Elizabeth','Morgan','Hannah','Alexis','Ashley','Samantha','Abigail','Olivia','Madison','Isabella','Grace','Brianna','Victoria','Natalie','Hailey','Alyssa','Savannah','Chloe','Claire','Rebecca','Leah','Audrey','Jasmine','Katherine','Mia','Makayla','Allison','Molly','Lily','Sydney','Kaylee','Ella','Maria','Kaitlyn','Brooke','Julia','Lillian','Taylor','Haley','Destiny','Katelyn','Faith','Mckenzie','Nicolas','Andrew','Alexander','Ryan','Daniel','Tyler','Matthew','Joseph','Logan','David','John','Brandon','Justin','Christopher','Jacob','Christian','Austin','Ethan','Benjamin','Noah','Samuel','Dylan','Nathan','Joshua','Gabriel','Jose','Anthony','William','Jonathan','Zachary','Kevin','Hunter','Caleb','Thomas','Kyle','Connor','Elijah','Jordan','Cameron','Jack','Evan','Luke','Eric','Isaac','Mason','Aaron','Luis','Charles','Eli','Juan','Adam','Sean','Nathaniel','Alex','Carlos','Ian','Bryan','Jesus','Owen','Lucas','Brian','Diego','Julian','Adrian','Gavin','Jaden','Chase','Cole','Antonio','Blake','Henry','Miguel','Sebastian','Timothy','Richard','Jake','Alejandro','Steven'])[floor(random()*223)+1] AS name_value,  -- Select a random name
    (ARRAY[true, false])[floor(random()*2)+1] AS bool_value,  -- Select a random boolean
    '2020-01-01'::date + (random() * 365 * 2)::integer AS date_value  -- Generate a random date
FROM generate_series(1, 10000000);


-- query by name, signup_date, active without index
EXPLAIN ANALYZE SELECT * FROM users WHERE name = 'Gary' AND signup_date BETWEEN '2021-09-01' AND '2021-12-31' AND active = true;

-- create inefficient index
CREATE INDEX idx_users_name_signup_date_active ON users (name, signup_date, active);

-- create efficient index (faster but it's not huge difference)
-- optimizer will pick this index over above index
CREATE INDEX idx_users_name_active_signup_date ON users (name, active, signup_date);
```

#### Inequality kills performance

```sql
-- create tables orders that has `status` field of type enum (pending, processing, shipped, delivered)
CREATE TYPE enum_order_status AS ENUM ('pending', 'processing', 'shipped', 'delivered');

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    status enum_order_status NOT NULL
);


-- seed data
INSERT INTO orders (status, user_id)
SELECT 'pending'::enum_order_status, 1
FROM generate_series(1, 10000000);
-- insert 10 records with status = 'processing'
-- insert 10 records with status = 'processing'
INSERT INTO orders (status, user_id)
SELECT 'processing'::enum_order_status, (ARRAY[2,3,4])[floor(random()*3)+1]
FROM generate_series(1, 10);
-- insert 10 records with status = 'shipped'
INSERT INTO orders (status, user_id)
SELECT 'shipped'::enum_order_status, (ARRAY[4,5,6])[floor(random()*3)+1]
FROM generate_series(1, 10);
-- insert 10 records with status = 'delivered'
INSERT INTO orders (status, user_id)
SELECT 'delivered'::enum_order_status, (ARRAY[6,7,8])[floor(random()*3)+1]
FROM generate_series(1, 10);


-- query
EXPLAIN ANALYZE SELECT * FROM orders WHERE status in ('processing', 'shipped', 'delivered');
EXPLAIN ANALYZE SELECT * FROM orders WHERE status = 'processing' OR status = 'shipped' OR status = 'delivered';
EXPLAIN ANALYZE SELECT * FROM orders WHERE status != 'pending';

-- create index
CREATE INDEX idx_orders_status ON orders (status);
EXPLAIN ANALYZE SELECT * FROM orders WHERE status != 'pending'; -- can't use index, seq scan
EXPLAIN ANALYZE SELECT * FROM orders WHERE status = 'processing' OR status = 'shipped' OR status = 'delivered'; -- can use index
EXPLAIN ANALYZE SELECT * FROM orders WHERE status in ('processing', 'shipped', 'delivered'); -- can use index

-- query on 2 fields
EXPLAIN ANALYZE SELECT * FROM orders WHERE status != 'pending' AND user_id = 4;
-- create index on (status, user_id) 
CREATE INDEX idx_orders_status_user_id ON orders (status, user_id); -- can use index, ~13ms
DROP INDEX idx_orders_status_user_id;
-- create index on (user_id, status)
CREATE INDEX idx_orders_user_id_status ON orders (user_id, status); -- can use index, ~0.3ms (faster)
-- why it's faster? because the first field in index is more selective
-- more selective: the number of rows that match a condition is small compared to the total number of rows in a table
```

#### Index on NULL values
```sql
-- create notificaitons table
CREATE TABLE notifications (
    notification_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    read_at TIMESTAMP WITH TIME ZONE
);

-- seed data 10.000.000 records with read_at is not null
INSERT INTO notifications (user_id, message, read_at)
SELECT (ARRAY[1,2,3,4,5,6,7,8,9,10])[floor(random()*10)+1], md5(random()::text), now()
FROM generate_series(1, 10000000);

-- seed 100 records with read_at is null for user_id in (1,2,3)
INSERT INTO notifications (user_id, message)
SELECT (ARRAY[1,2,3])[floor(random()*3)+1], md5(random()::text)
FROM generate_series(1, 100);

-- seed 30 records with read_at is null for user_id = 15
INSERT INTO notifications (user_id, message)
SELECT 15, md5(random()::text)
FROM generate_series(1, 30);

-- seed 30 records with read_at not null for user_id = 20
INSERT INTO notifications (user_id, message, read_at)
SELECT 20, md5(random()::text), now()
FROM generate_series(1, 30);


-- query
explain analyze select * from notifications where read_at is null and user_id = 1; -- >300ms
------------------
-- CASE: create index on (user_id, read_at)
CREATE INDEX idx_notifications_user_id_read_at ON notifications (user_id, read_at);
explain analyze select * from notifications where read_at is null and user_id = 1; -- 1ms
-- check index size
SELECT pg_size_pretty(pg_relation_size('idx_notifications_user_id_read_at')); -- 67MB

------------------
-- CASE: create index on (user_id, read_at) where read_at is null (as we only care about unread notifications)
CREATE INDEX idx_notifications_user_id_read_at_null ON notifications (user_id, read_at) WHERE read_at IS NULL;
-- check index size
SELECT pg_size_pretty(pg_relation_size('idx_notifications_user_id_read_at_null')); -- 16KB (only index unread notifications)
explain analyze select * from notifications where read_at is null and user_id = 1; -- <1ms (this index is chosen by optimizer)

------------------
-- CASE: create index on (read_at, user_id)
CREATE INDEX idx_notifications_read_at_user_id ON notifications (read_at, user_id);
explain analyze select * from notifications where read_at is null and user_id = 1; -- this index is chosen by optimizer

------------------
-- CASE: create index on (read_at, user_id) where read_at is null
CREATE INDEX idx_notifications_read_at_user_id_null ON notifications (read_at, user_id) WHERE read_at IS NULL;
explain analyze select * from notifications where read_at is null and user_id = 1; -- this index is chosen by optimizer
```

#### Pattern matching
```sql
-- create table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    country TEXT NOT NULL
);

-- seed junk data 10.000.000 records
INSERT INTO users (first_name, last_name, country)
SELECT 
    md5(random()::text), 
    md5(random()::text), 
    (ARRAY['ABC','DEF','XXX','MNP','QQQ'])[floor(random()*5)+1]
FROM generate_series(1, 10000000);

-- seed deterministic data
INSERT INTO users (first_name, last_name, country) VALUES 
    ('John', 'Doe', 'US'),
    ('John', 'Doe', 'UK'),
    ('John', 'Doe', 'FR'),
    ('John', 'Doe', 'IT'),
    ('Steve', 'Jobs', 'US'),
    ('Steve', 'Jobs', 'UK'),
    ('Steve', 'Jobs', 'FR'),
    ('Steve', 'Jobs', 'IT'),
    ('Steve', 'Wozniak', 'US'),
    ('Steve', 'Wozniak', 'UK'),
    ('Steve', 'Wozniak', 'FR'),
    ('Steve', 'Wozniak', 'IT');

-- query with LIKE
EXPLAIN ANALYZE SELECT * FROM users WHERE first_name LIKE 'Joh%'; -- 500ms

-- create index on first_name
CREATE INDEX idx_users_first_name ON users (first_name);
EXPLAIN ANALYZE SELECT * FROM users WHERE first_name LIKE 'Joh%'; -- 500ms, no index used
-- create index on first_name and country
CREATE INDEX idx_users_first_name_country ON users (first_name, country);
EXPLAIN ANALYZE SELECT * FROM users WHERE first_name LIKE 'Joh%'; -- 500ms, no index used
EXPLAIN ANALYZE SELECT * FROM users WHERE first_name LIKE 'Joh%' AND country = 'US';
-- create index on country and first_name
CREATE INDEX idx_users_country_first_name ON users (country, first_name);
EXPLAIN ANALYZE SELECT * FROM users WHERE first_name LIKE 'Joh%'; -- 500ms, no index used
EXPLAIN ANALYZE SELECT * FROM users WHERE first_name LIKE 'Joh%' AND country = 'US'; -- 0.2ms, index used `idx_users_country_first_name` 
```

#### Sorting values
```sql
-- create table github_issues with status, severity, num_of_comments
CREATE TABLE github_issues (
    issue_id SERIAL PRIMARY KEY,
    status TEXT NOT NULL,
    severity TEXT NOT NULL,
    num_of_comments INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- seed data 10.000.000 records
INSERT INTO github_issues (status, severity, num_of_comments, created_at)
SELECT 
    (ARRAY['open','closed','assigned','resolved'])[floor(random()*4)+1], 
    (ARRAY[1,2,3,4,5,6,7,8,9,10])[floor(random()*10)+1], 
    (ARRAY[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50])[floor(random()*51)+1], 
    now() - (random() * 365 * 2)::integer * '1 day'::interval
FROM generate_series(1, 10000000);

-- query
explain analyze select * from github_issues 
where status = 'open' 
order by severity desc, num_of_comments desc, created_at desc limit 20; -- extremely fast
-- index
CREATE INDEX idx_github_issues_status_severity_num_of_comments_created_at ON github_issues (status, severity, num_of_comments, created_at);
-- experiment: drop above index, create index on (status)
CREATE INDEX idx_github_issues_status ON github_issues (status); -- very slow
-- experiment: drop above index, create index on (status, severity) 
CREATE INDEX idx_github_issues_status_severity ON github_issues (status, severity); -- faster
-- experiment: drop above index, create index on (status, severity, num_of_comments)
CREATE INDEX idx_github_issues_status_severity_num_of_comments ON github_issues (status, severity, num_of_comments);
```

#### Demo `work_mem`: Adjust such that sort can be done in memory
```sql
-- sample table
CREATE TABLE test_data (
    id serial PRIMARY KEY,
    data_value int,
    created_at timestamp default current_timestamp
);

-- seed data 10.000.000 records
INSERT INTO test_data (data_value, created_at)
SELECT (random()*100000)::int, 
       timestamp '2021-01-01 00:00:00' + interval '1 minute' * (random()*525600)::int
FROM generate_series(1,1000000);
-- 50 bytes
-- 1MB * 1024 * 1000 bytes / 50 bytes = 20480 rows

-- show current work_mem
SHOW work_mem;
--  work_mem 
-- ----------
--  4MB
-- (1 row)

-- reduce work_mem to 1MB
SET work_mem = '1MB';

-- query that exceeds work_mem
EXPLAIN ANALYZE SELECT data_value, created_at FROM test_data WHERE created_at > '2020-01-01' ORDER BY data_value DESC LIMIT 40000;
-- Sort Method: external merge  Disk: 8552kB
-- Worker 0:  Sort Method: external merge  Disk: 8472kB
-- Worker 1:  Sort Method: external merge  Disk: 8504kB

-- query that doesn't exceed work_mem
SELECT count(*) FROM test_data WHERE created_at > '2021-12-30';
EXPLAIN ANALYZE SELECT data_value, created_at FROM test_data WHERE created_at > '2021-12-30' ORDER BY data_value DESC; -- ~5000 rows
-- Sort Method: quicksort  Memory: 219kB
-- Worker 0:  Sort Method: quicksort  Memory: 111kB
-- Worker 1:  Sort Method: quicksort  Memory: 117kB

SELECT count(*) FROM test_data WHERE created_at > '2021-09-30'; -- ~25000 rows
EXPLAIN ANALYZE SELECT data_value, created_at FROM test_data WHERE created_at > '2021-09-30' ORDER BY data_value DESC; -- ~25000 rows
-- Sort Method: external merge  Disk: 2224kB
-- Worker 0:  Sort Method: external merge  Disk: 2136kB
-- Worker 1:  Sort Method: external merge  Disk: 2176kB

-- increase work_mem to 2MB
SET work_mem = '4MB';
EXPLAIN ANALYZE SELECT data_value, created_at FROM test_data WHERE created_at > '2021-09-30' ORDER BY data_value DESC; -- ~25000 rows
-- Sort Method: quicksort  Memory: 6061kB
-- Worker 0:  Sort Method: quicksort  Memory: 5993kB
-- Worker 1:  Sort Method: quicksort  Memory: 6011kB
```


#### Demo index for `GROUP BY` on 1 field
```sql
-- create table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    country TEXT NOT NULL
);

-- seed data 10.000.000 records with country in 240 countries
INSERT INTO users (first_name, last_name, country)
SELECT 
    md5(random()::text), 
    md5(random()::text), 
    (ARRAY['AF','AL','DZ','AS','AD','AO','AI','AQ','AG','AR','AM','AW','AU','AT','AZ','BS','BH','BD','BB','BY','BE','BZ','BJ','BM','BT','BO','BA','BW','BR','IO','VG','BN','BG','BF','BI','KH','CM','CA','CV','KY','CF','TD','CL','CN','CX','CC','CO','KM','CK','CR','HR','CU','CW','CY','CZ','CD','DK','DJ','DM','DO','TL','EC','EG','SV','GQ','ER','EE','ET','FK','FO','FJ','FI','FR','PF','GA','GM','GE','DE','GH','GI','GR','GL','GD','GU','GT','GG','GN','GW','GY','HT','HN','HK','HU','IS','IN','ID','IR','IQ','IE','IM','IL','IT','CI','JM','JP','JE','JO','KZ','KE','KI','XK','KW','KG','LA','LV','LB','LS','LR','LY','LI','LT','LU','MO','MK','MG','MW','MY','MV','ML','MT','MH','MR','MU','YT','MX','FM','MD','MC','MN','ME','MS','MA','MZ','MM','NA','NR','NP','NL','AN','NC','NZ','NI','NE','NG','NU','KP','MP','NO','OM','PK','PW','PS','PA','PG','PY','PE','PH','PN','PL','PT','PR','QA','CG','RE','RO','RU','RW','BL','SH','KN','LC','MF','PM','VC','WS','SM','ST','SA','SN','RS','SC','SL','SG','SX','SK','SI','SB','SO','ZA','KR','SS','ES','LK','SD','SR','SJ','SZ','SE','CH','SY','TW','TJ','TZ','TH','TG','TK','TO','TT','TN','TR','TM','TC','TV','VI','UG','UA','AE','GB','US','UY','UZ','VU','VA','VE','VN','WF','EH','YE','ZM','ZW'])[floor(random() * 240) + 1]
FROM generate_series(1, 10000000);

-- query
EXPLAIN ANALYZE SELECT country, count(*) FROM users GROUP BY country;

-- create index
CREATE INDEX idx_users_country ON users (country);

-- query
EXPLAIN ANALYZE SELECT country, count(*) FROM users GROUP BY country;
EXPLAIN ANALYZE SELECT country, count(*) FROM users WHERE country IN ('VN', 'SG') GROUP BY country;

-- test=# select query, total_exec_time, mean_exec_time, calls, rows from pg_stat_statements where query ilike '%from users%';
--                                                query                                                | total_exec_time |  mean_exec_time   | calls | rows 
-- ----------------------------------------------------------------------------------------------------+-----------------+-------------------+-------+------
--  SELECT country, count(*) FROM users GROUP BY country                                               |       3616.8505 |        904.212625 |     4 |  960
--  SELECT country, count(*) FROM users WHERE country IN ($1, $2) GROUP BY country                     |       5701.8596 | 633.5399555555556 |     9 |   18

-- reset pg_stat_statements
SELECT pg_stat_statements_reset();

-- test=# select query, total_exec_time, mean_exec_time, calls, rows from pg_stat_statements where query ilike '%from users%';
--                                      query                                      |  total_exec_time   |   mean_exec_time   | calls | rows 
-- --------------------------------------------------------------------------------+--------------------+--------------------+-------+------
--  SELECT country, count(*) FROM users GROUP BY country                           | 2371.7812000000004 |  592.9453000000001 |     4 |  960
--  SELECT country, count(*) FROM users WHERE country IN ($1, $2) GROUP BY country |           154.7948 | 14.072254545454545 |    11 |   22

```

## DEMO index for `GROUP BY` on 2 fields
```sql
-- DEMO index for `GROUP BY` on 2 fields
-- create table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    country TEXT NOT NULL,
    star_sign TEXT NOT NULL
);

-- seed data 10.000.000 records with country in 240 countries
INSERT INTO users (first_name, last_name, country, star_sign)
SELECT 
    md5(random()::text), 
    md5(random()::text),
    (ARRAY['AF','AL','DZ','AS','AD','AO','AI','AQ','AG','AR','AM','AW','AU','AT','AZ','BS','BH','BD','BB','BY','BE','BZ','BJ','BM','BT','BO','BA','BW','BR','IO','VG','BN','BG','BF','BI','KH','CM','CA','CV','KY','CF','TD','CL','CN','CX','CC','CO','KM','CK','CR','HR','CU','CW','CY','CZ','CD','DK','DJ','DM','DO','TL','EC','EG','SV','GQ','ER','EE','ET','FK','FO','FJ','FI','FR','PF','GA','GM','GE','DE','GH','GI','GR','GL','GD','GU','GT','GG','GN','GW','GY','HT','HN','HK','HU','IS','IN','ID','IR','IQ','IE','IM','IL','IT','CI','JM','JP','JE','JO','KZ','KE','KI','XK','KW','KG','LA','LV','LB','LS','LR','LY','LI','LT','LU','MO','MK','MG','MW','MY','MV','ML','MT','MH','MR','MU','YT','MX','FM','MD','MC','MN','ME','MS','MA','MZ','MM','NA','NR','NP','NL','AN','NC','NZ','NI','NE','NG','NU','KP','MP','NO','OM','PK','PW','PS','PA','PG','PY','PE','PH','PN','PL','PT','PR','QA','CG','RE','RO','RU','RW','BL','SH','KN','LC','MF','PM','VC','WS','SM','ST','SA','SN','RS','SC','SL','SG','SX','SK','SI','SB','SO','ZA','KR','SS','ES','LK','SD','SR','SJ','SZ','SE','CH','SY','TW','TJ','TZ','TH','TG','TK','TO','TT','TN','TR','TM','TC','TV','VI','UG','UA','AE','GB','US','UY','UZ','VU','VA','VE','VN','WF','EH','YE','ZM','ZW'])[floor(random() * 240) + 1],
    (ARRAY['Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo', 'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces'])[floor(random() * 12) + 1]
FROM generate_series(1, 10000000);

-- query without index
SELECT country, star_sign, count(*) FROM users GROUP BY country, star_sign;
--                                      query                                      |  total_exec_time   |   mean_exec_time   | calls | rows  
-- --------------------------------------------------------------------------------+--------------------+--------------------+-------+-------
--  SELECT country, star_sign, count(*) FROM users GROUP BY country, star_sign     |  5939.052599999999 | 1187.8105199999998 |     5 | 14400

-- create index
CREATE INDEX idx_users_country_star_sign ON users (country, star_sign);
SELECT country, star_sign, count(*) FROM users GROUP BY country, star_sign;
--                                    query                                    |  total_exec_time  |  mean_exec_time   | calls | rows  
-- ----------------------------------------------------------------------------+-------------------+-------------------+-------+-------
--  SELECT country, star_sign, count(*) FROM users GROUP BY country, star_sign | 4343.040800000001 | 723.8401333333334 |     6 | 17280

-- DEMO: GROUP BY 2 fields with WHERE filter
-- insert deterministic data: 100 records with country in ('VN', 'SG', 'MY') and star_sign in ('Aries', 'Taurus', 'Gemini') and first_name in ('John', 'Steve', 'Bill')
INSERT INTO users (first_name, last_name, country, star_sign)
SELECT 
    (ARRAY['John', 'Steve', 'Bill'])[floor(random() * 3) + 1],
    md5(random()::text),
    (ARRAY['VN', 'SG', 'MY'])[floor(random() * 3) + 1],
    (ARRAY['Aries', 'Taurus', 'Gemini'])[floor(random() * 3) + 1]
FROM generate_series(1, 100);

-- query without index
SELECT country, star_sign, count(*) FROM users 
WHERE first_name IN ('John', 'Steve', 'Bill') GROUP BY country, star_sign;

-- create index
CREATE INDEX idx_users_first_name_country_star_sign ON users (first_name, country, star_sign);

-- DEMO: GROUP BY 2 fields with WHERE range filter
-- create table users (country, star_sign, balance)
create table users (
    user_id serial primary key,
    country text not null,
    star_sign text not null,
    balance integer not null
);

-- seed data 10.000.000 records with country in 240 countries, 12 star_signs, balance in range [0, 1000]
INSERT INTO users (country, star_sign, balance)
SELECT
    (ARRAY['AF','AL','DZ','AS','AD','AO','AI','AQ','AG','AR','AM','AW','AU','AT','AZ','BS','BH','BD','BB','BY','BE','BZ','BJ','BM','BT','BO','BA','BW','BR','IO','VG','BN','BG','BF','BI','KH','CM','CA','CV','KY','CF','TD','CL','CN','CX','CC','CO','KM','CK','CR','HR','CU','CW','CY','CZ','CD','DK','DJ','DM','DO','TL','EC','EG','SV','GQ','ER','EE','ET','FK','FO','FJ','FI','FR','PF','GA','GM','GE','DE','GH','GI','GR','GL','GD','GU','GT','GG','GN','GW','GY','HT','HN','HK','HU','IS','IN','ID','IR','IQ','IE','IM','IL','IT','CI','JM','JP','JE','JO','KZ','KE','KI','XK','KW','KG','LA','LV','LB','LS','LR','LY','LI','LT','LU','MO','MK','MG','MW','MY','MV','ML','MT','MH','MR','MU','YT','MX','FM','MD','MC','MN','ME','MS','MA','MZ','MM','NA','NR','NP','NL','AN','NC','NZ','NI','NE','NG','NU','KP','MP','NO','OM','PK','PW','PS','PA','PG','PY','PE','PH','PN','PL','PT','PR','QA','CG','RE','RO','RU','RW','BL','SH','KN','LC','MF','PM','VC','WS','SM','ST','SA','SN','RS','SC','SL','SG','SX','SK','SI','SB','SO','ZA','KR','SS','ES','LK','SD','SR','SJ','SZ','SE','CH','SY','TW','TJ','TZ','TH','TG','TK','TO','TT','TN','TR','TM','TC','TV','VI','UG','UA','AE','GB','US','UY','UZ','VU','VA','VE','VN','WF','EH','YE','ZM','ZW'])[floor(random() * 240) + 1],
    (ARRAY['Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo', 'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces'])[floor(random() * 12) + 1],
    floor(random() * 1000)
FROM generate_series(1, 10000000);

-- query without index
SELECT country, star_sign, count(*) FROM users 
WHERE balance BETWEEN 0 AND 300 GROUP BY country, star_sign;
SELECT country, star_sign, count(*) FROM users 
WHERE balance BETWEEN 301 AND 700 GROUP BY country, star_sign;
SELECT country, star_sign, count(*) FROM users 
WHERE balance BETWEEN 701 AND 1000 GROUP BY country, star_sign;
-- => ~ 783ms

-- create index on (balance, country, star_sign)
CREATE INDEX idx_users_balance_country_star_sign ON users (balance, country, star_sign);
-- => ~ 611 ms (did enhance performance but not much)

-- experiment: create index on (country, star_sign, balance)
DROP INDEX idx_users_balance_country_star_sign;
CREATE INDEX idx_users_country_star_sign_balance ON users (country, star_sign, balance);
-- => ~ 835ms (not good)

-- experiment: calculated column `balance_range` and create index on (country, star_sign, balance_range)
ALTER TABLE users 
ADD COLUMN balance_range INTEGER GENERATED ALWAYS AS (
    CASE WHEN balance BETWEEN 0 AND 300 THEN 1 
    WHEN balance BETWEEN 301 AND 700 THEN 2 
    WHEN balance BETWEEN 701 AND 1000 THEN 3 END
) STORED;
CREATE INDEX idx_users_balance_range_country_star_sign ON users (balance_range, country, star_sign);

-- using avg to calculate average balance of each group (country, star_sign)

-- query without index
SELECT country, star_sign, avg(balance) FROM users 
WHERE balance BETWEEN 0 AND 300 GROUP BY country, star_sign;

-- create index on (country, star_sign, balance)
CREATE INDEX idx_users_country_star_sign_balance ON users (country, star_sign, balance);

```

#### Indexes applied for JOIN
```sql
-- create table employees, departments
CREATE TABLE departments (
    department_id SERIAL PRIMARY KEY,
    country TEXT NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE employees (
    employee_id SERIAL PRIMARY KEY,
    salary DECIMAL NOT NULL,
    department_id INTEGER NOT NULL REFERENCES departments(department_id)
);


-- seed data 1.000 departments of 50 countries
INSERT INTO departments (country, name)
SELECT 
    (ARRAY[
  'AF', 'AL', 'DZ', 'AS', 'AD', 'AO',
  'AI', 'AQ', 'AG', 'AR', 'AM', 'AW',
  'AU', 'AT', 'AZ', 'BS', 'BH', 'BD',
  'BB', 'BY', 'BE', 'BZ', 'BJ', 'BM',
  'BT', 'BO', 'BA', 'BW', 'BR', 'IO',
  'VG', 'BN', 'BG', 'BF', 'BI', 'KH',
  'CM', 'CA', 'CV', 'KY', 'CF', 'TD',
  'CL', 'CN', 'CX', 'CC', 'CO', 'KM',
  'CK', 'CR'
])[floor(random() * 50) + 1],
    md5(random()::text)
FROM generate_series(1, 1000);

-- insert 1.000.000 employees with salary in range [1000, 2000] and department_id in range [1, 1000]
INSERT INTO employees (salary, department_id)
SELECT 
    (random() * 1000 + 1000)::integer,
    floor(random() * 1000) + 1
FROM generate_series(1, 1000000);

-- query employees with salary > 1500 and department.country = 'AU'
SELECT employees.employee_id, employees.salary, departments.country, departments.name
FROM employees
INNER JOIN departments ON employees.department_id = departments.department_id
WHERE employees.salary > 1500 AND departments.country = 'AU';

-- create index on (salary)
CREATE INDEX idx_employees_salary ON employees (salary);
-- create index on (country, department_id)
CREATE INDEX idx_departments_country_department_id ON departments (country, department_id);
-- still slow
-- game changer:
CREATE INDEX idx_employees_department_id_salary ON employees (department_id, salary);
```

### Index applied for subquery
```sql
-- given this query
SELECT *
FROM products
WHERE remaining > 500 AND category_id = (
  SELECT category_id
  FROM categories
  WHERE type = 'game' AND name = '83f23e1bdccd70ad8b691b24b0e83f59'
)

-- create table products, categories
CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    remaining INTEGER NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories(category_id)
);

-- seed data 1000 categories of 10 types
INSERT INTO categories (type, name)
SELECT 
    (ARRAY['book', 'movie', 'music', 'game', 'food', 'drink', 'clothes', 'shoes', 'accessories', 'electronics'])[floor(random() * 10) + 1],
    md5(random()::text)
FROM generate_series(1, 1000);

-- seed data 1.000.000 products with remaining in range [0, 1000] and category_id in range [1, 1000]
INSERT INTO products (name, remaining, category_id)
SELECT 
    md5(random()::text),
    floor(random() * 1000),
    floor(random() * 1000) + 1
FROM generate_series(1, 1000000);

-- no index at all: 50ms

-- create index on categories (type, name)
CREATE INDEX idx_categories_type_name ON categories (type, name); -- 42ms

-- create index on categories (name, type, category_id)
CREATE INDEX idx_categories_name_type_category_id ON categories (name, type, category_id); -- chosen by optimizer

-- create index on products (category_id, remaining)
CREATE INDEX idx_products_category_id_remaining ON products (category_id, remaining); -- extremely fast 1ms

-- experiment: create index on products (remaining, category_id)
DROP INDEX idx_products_category_id_remaining;
CREATE INDEX idx_products_remaining_category_id ON products (remaining, category_id); -- 20ms

```

#### Index applied for dependent subquery
```sql
-- given this query: nested loop
-- for each matching product, run a loop to check if there is any sale in the last 6 months in sales table
SELECT *
FROM products
WHERE remaining = 0 AND EXISTS (
  SELECT *
  FROM sales
  WHERE created_at >= '2023-06-01' AND product_id = products.product_id
);

-- create sales table
CREATE TABLE sales (
    sale_id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(product_id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- seed data 1.000.000 sales with created_at in range ['2023-01-01', '2024-01-01'] and product_id in range [1, 1000000]
INSERT INTO sales (product_id, created_at)
SELECT 
    floor(random() * 1000000) + 1,
    timestamp '2023-01-01' + interval '1 minute' * (random() * 525600)
FROM generate_series(1, 1000000);

-- without index: 160ms
-- create index on sales products(remaining)
CREATE INDEX idx_products_remaining ON products (remaining); -- index not used

-- experiment: create index on sales (remaining, product_id)
CREATE INDEX idx_products_remaining_product_id ON products (remaining, product_id); -- index not used

-- experiment: create index on sales (product_id, created_at)
CREATE INDEX idx_sales_product_id_created_at ON sales (product_id, created_at); -- 3ms
-- 2 indexes are used: `idx_products_remaining` and `idx_sales_product_id_created_at`

-- experiment: drop index idx_products_remaining
DROP INDEX idx_products_remaining;
-- fallback to idx_products_remaining_product_id

-- experiment: drop index idx_products_remaining_product_id
-- becomes slow 60ms
```