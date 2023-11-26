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