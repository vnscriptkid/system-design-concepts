## https://chat.openai.com/c/6f5ca7dd-5d5c-4eab-a889-df78c0947aa5

## Summary
- Isolation is a key component of the ACID properties of relational databases, ensuring that concurrent transactions do not interfere with each other.
- Isolation levels dictate the visibility of changes made by one transaction to other concurrent transactions and control phenomena like dirty reads, non-repeatable reads, and phantom reads.
- The Serializable isolation level offers the highest level of isolation by making every read a locking read, which means while one transaction is reading data, other transactions must wait, thus preventing concurrent access issues.
- It is important to note that storage engines can alter the implementation of isolation levels, so it's essential to read the documentation before making changes to the database configuration.


## Exercises
1. Setting Up the Environment: Create a simple database schema with a table for accounts including fields for account ID, user ID, and balance. Populate it with sample data.
```sql
CREATE TABLE accounts (
  id SERIAL NOT NULL PRIMARY KEY,
  user_id INT NOT NULL,
  balance DECIMAL(10,2) NOT NULL
);

INSERT INTO accounts (id, user_id, balance) VALUES (1, 1, 1000.00);
```

2. Read Uncommitted Demonstration: Start two transactions. In the first, update the balance of an account but do not commit. In the second transaction, read the updated balance before the first transaction commits. Discuss the implications of a dirty read.
```sql
-- TERM1
BEGIN;
UPDATE accounts SET balance = balance + 100 WHERE id = 1;
```

```sql
-- TERM2
BEGIN;
-- current isolation level
SHOW TRANSACTION ISOLATION LEVEL;
--> read committed
-- set isolation level in postgres
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
SHOW TRANSACTION ISOLATION LEVEL;
--> read uncommitted

SELECT balance FROM accounts WHERE id = 1;
--> 1100.00 (dirty read, in mysql)
--> 1000.00 (read uncommitted ~ read committed, in postgres)
```

3. Read Committed Exercise: Repeat the previous exercise but set the isolation level to Read Committed. Observe that the second transaction cannot see the uncommitted balance update.

4. Repeatable Reads Scenario: Set the isolation level to Repeatable Reads. Perform a select statement in one transaction, update the data in another transaction, and then perform the select statement again in the first transaction. Discuss why the results are the same despite the update.
```sql
-- TERM1
BEGIN;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SELECT balance FROM accounts WHERE id = 1;
--> 1000.00

--- TERM2
BEGIN;
UPDATE accounts SET balance = balance + 100 WHERE id = 1;
COMMIT;

-- TERM1
SELECT balance FROM accounts WHERE id = 1;
--> 1000.00

-- All scenarios:
-- TERM1 begins, TERM2 begins, TERM1 selects, TERM2 updates, TERM2 commits, TERM1 selects
-- TERM1 begins, TERM2 begins, TERM2 updates, TERM2 commits, TERM1 selects
-- TERM1 begins, TERM2 begins, TERM2 updates, TERM1 selects, TERM2 commits, TERM1 selects
```

5. Phantom Read Exploration: Using the Repeatable Reads isolation level, insert new rows in one transaction and then perform a range SELECT in another transaction, twice. Observe the effect of phantom reads.
```sql
-- TERM1
BEGIN;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SELECT * FROM accounts WHERE balance > 500;

-- TERM2
BEGIN;
INSERT INTO accounts (id, user_id, balance) VALUES (2, 2, 1000.00);
COMMIT;

-- TERM1
SELECT * FROM accounts WHERE balance > 500;
```

6. Serializable Transactions: Set the isolation level to Serializable and attempt concurrent transactions that modify the same data. Observe how one transaction blocks the other and discuss the benefits and drawbacks.
```sql
-- TERM1
BEGIN;
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
UPDATE accounts SET balance = balance + 100 WHERE id = 1;

-- TERM2
BEGIN;
SELECT * FROM accounts WHERE id = 2; -- can still read
SELECT * FROM accounts WHERE id = 1; -- can still read
UPDATE accounts SET balance = balance + 100 WHERE id = 2; -- can still update other rows
UPDATE accounts SET balance = balance + 100 WHERE id = 1; -- write blocked (row level lock)
COMMIT;

-- TERM1  
COMMIT; -- after this, TERM2 can update as lock is released
```

7. Isolation Level Impact on Performance: Run a series of SELECT and UPDATE operations under different isolation levels and measure the execution time. Discuss how stricter isolation levels can affect database performance.

8. Lost Updates: Create a scenario where two transactions read the same data and then update it based on the read. Discuss how the last commit can overwrite changes made by the first commit and how to prevent this.
```sql
-- Begin two transactions under the Read Committed isolation level
-- TERM1
BEGIN;
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
SELECT balance FROM accounts WHERE id = 1;

-- TERM2
BEGIN;
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
SELECT balance FROM accounts WHERE id = 1;

-- TERM1
UPDATE accounts SET balance = balance + 100 WHERE id = 1;
COMMIT;

-- TERM2
UPDATE accounts SET balance = balance + 100 WHERE id = 1; -- overwrites TERM1's update
COMMIT;

-- TERM1
SELECT balance FROM accounts WHERE id = 1;
--> 1200.00

-- TERM2
SELECT balance FROM accounts WHERE id = 1;
--> 1200.00
```

9. Isolation Levels and Locking: Show how different isolation levels correspond to different types of locking (row-level, table-level, etc.) by observing database behavior when multiple transactions access the same data.

10. Implementing Snapshot Isolation: If supported by the database system, demonstrate how Snapshot Isolation can be used to provide a consistent view of the data without locking, and discuss how it differs from the other levels.
