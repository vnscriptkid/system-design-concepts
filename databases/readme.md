## Topics
1. Relational Databases
2. Database Isolation Levels
3. Scaling Databases
4. Sharding and Paritioning
5. Non-relational Databases
6. When to use which database

## Deeper
1. List all of database constraints?
- Primary Key
- Foreign Key
- Unique
- Not Null
- Check

2. List all types of triggers?
- Before
- After
- Instead of (TODO: research this)

3. Trigger inside transaction? If trigger fails, does transaction fail?
- Yes, trigger is part of transaction

4. Usecases of triggers?
- Logging. For example, if you want to log all changes to a table, you can create a trigger that fires after an update, insert, or delete statement.
- Auditing. For example, if you want to keep track of all users who are making changes to a table, you can create a trigger that fires before an update, insert, or delete statement. The trigger can then record the user's name, the action, and the time in an audit table.
- Enforcing business rules: For example, if you want to enforce a rule that states that all customers must have a credit limit, you can create a trigger that fires before an insert or update statement and cancels the insert or update if the credit limit is null.
- Data integrity: For example, if you want to ensure that all orders have at least one line item, you can create a trigger that fires before an insert or update statement and cancels the insert or update if the order does not have any line items.
- Replication: For example, if you want to replicate all changes to a table in another database, you can create a trigger that fires after an update, insert, or delete statement and records the changes in a separate table.
- Synchronization: For example, if you want to synchronize the data in two tables, you can create a trigger that fires after an update, insert, or delete statement and updates the data in the other table. An practical use case of this is to synchronize the data between two databases that are located in different locations.

5. Why not to use triggers?
- Performance: Triggers can slow down the database because they are executed on every row affected by the triggering statement.
- Complexity: Triggers are hard to debug and maintain because they are not part of the application code.
- Portability: Triggers are not portable because they are written in the native language of the database server.

6. Rollback
- Explicit rollback in try/catch block: triggered by application
    - Fails business rule validation
    - Failure in external system
    ```go
    tx, err := db.Begin()

    err := tx.Exec("INSERT INTO metadata(id, file_type) VALUES (?, ?)", "123", "png")
    if err != nil {
        tx.Rollback()
        return err
    }

    err := uploadToS3()
    if err != nil {
        tx.Rollback()
        return err
    }

    tx.Commit()
    ```
- Implicit rollback when error occurs (e.g. constraint violation): triggered by database

7. ACID
- Atomicity: All or nothing
- Consistency: Data is valid/consistent/accurate
- Isolation: How much is visible to other transactions
- Durability: Data is persisted

9. Durability
- Write-ahead logging (WAL): 
```json
{
  "WAL_Record": {
    "type": "update",
    "transaction_id": 12345,
    "relation": "accounts",
    "data_block": 4321,
    "before_image": {
      "account_id": 1,
      "balance": 1000
    },
    "after_image": {
      "account_id": 1,
      "balance": 900
    }
  }
}
```
- What COMMIT does internally:
    - Writes a Commit Record to the WAL (WAL now stored in memory)
    - Flushes the WAL Records from RAM to disk (after this, the transaction is considered durable)
        - If system fails after this, the WAL can be replayed to recover the database
    - Database confirms to the client that the transaction is committed
    - Actual data is written to disk. 
        - Some databases write to disk immediately
        - Some databases write to disk periodically. This is called a checkpoint.


10. Implicit transactions/autocommit
- If you don't explicitly start a transaction, every statement is its own transaction
```sql
-- This is a transaction
INSERT INTO accounts (account_id, balance) VALUES (1, 1000);
-- Batch insert is a transaction
INSERT INTO accounts (account_id, balance) VALUES 
    (2, 2000), (3, 3000);
```

11. ON DELETE:
- ON DELETE CASCADE:
    - Brutal. Why? Maintain correct data integrity
- ON DELETE SET NULL: 
    - Does not maintain data integrity. Leads to orphaned records (no parent)
- ON DELETE RESTRICT (default): 
    - Strict. Why? Maintain correct data integrity

12. Explain lost update problem
- Changes made by one transaction are overwritten by another transaction
- Happens in isolation level READ COMMITTED, READ UNCOMMITTED
- Does not happen in isolation level REPEATABLE READ, SERIALIZABLE

13. REPEATABLE READ behaviors across databases
- MySQL: 
    - Locks rows that are read
    - Locks rows that are written
- PostgreSQL:
    - Locks rows that are written
    - Does not lock rows that are read
        - Why? PostgreSQL uses MVCC (Multi-Version Concurrency Control). This means that when a row is updated, the old version of the row is kept in the database. This allows other transactions to read the old version of the row. This is why PostgreSQL does not need to lock rows that are read.

14. READ UNCOMMITTED behaviors across databases
- MySQL:
    - susceptible to dirty reads
- PostgreSQL:
    - Just like READ COMMITTED, does not allow dirty reads
    - Lock
        - Does not lock rows that are read
        - Does not lock rows that are written

15. READ COMMITTED behaviors across databases
- MySQL:
- PostgreSQL:
    - Does not allow dirty reads
    - Does not lock rows that are read
    - Does not lock rows that are written 
        - susceptible to lost update problem, last write wins
        - Fix:
            - Locking: `SELECT ... FOR UPDATE`
            - Optimistic locking: `UPDATE ... WHERE balance = 1000`

16. Implicit locking in PostgreSQL across isolation levels
- READ COMMITTED:
    - Does not lock rows that are read
    - Does not lock rows that are written
- REPEATABLE READ:
    - Does not lock rows that are read
    - Locks rows that are written
- SERIALIZABLE: (same as REPEATABLE READ)
    - Does not lock rows that are read
    - Locks rows that are written

17. SERIALIZABLE in PostgreSQL
- What does word "serializable" mean?
    - Transactions are executed one after another (serial order): T1 -> T2 || T2 -> T1
- How Postgres check for serializability error?
    - Check dependencies between transactions
- Consistency guarantee by:
    - Locking: immediate consistency for writes (manage write conflicts)
    - SSI: broad consistency by ensuring overall serializability
        - Checks for serializability errors
        - If error, abort transaction
        - If no error, commit transaction
- Serialization error usecase 1:
    - T1 and T2 are not serializable because T1 reads the sum of all balances before T2 commits. If T1 were to read the sum of all balances after T2 commits, then T1 and T2 would be serializable.

```sql

CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  balance INT
);

TRUNCATE TABLE accounts;
-- Insert initial data
INSERT INTO accounts (id, balance) VALUES (1, 100);
INSERT INTO accounts (id, balance) VALUES (2, 200);
```

| Time Order | Operation in Transaction 1 (T1)                     | Operation in Transaction 2 (T2)              |
|------------|-----------------------------------------------------|----------------------------------------------|
| 1          | `BEGIN;`                                            |                                              |
| 2          | `SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;`     |                                              |
| 3          | `SELECT SUM(balance) FROM accounts;`                |                                              |
| 4          |                                                     | `BEGIN;`                                     |
| 5          |                                                     | `SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;` |
| 6          |                                                     | `UPDATE accounts SET balance = 250 WHERE id = 1;` |
| 7          |                                                     | `COMMIT;`                                    |
| 8          | `INSERT INTO accounts (balance) VALUES (300);`      |                                              |
| 9          | `COMMIT;`                                           |                                              |

ERROR:  could not serialize access due to read/write dependencies among transactions
DETAIL:  Reason code: Canceled on identification as a pivot, during write.
HINT:  The transaction might succeed if retried.

- Line 8: Postgres considers the possibility of dependency: INSERT might be influenced by results of SELECT,
which has been invalidated by UPDATE in T2

18. Real world usecase of REPEATABLE READ
- Need calculation based on multiple rows and make decision based on that calculation
    - Calculation remains consistent throughout transaction

19. READ COMMITTED with LOCKING
```sql
CREATE TABLE events (
    event_id SERIAL PRIMARY KEY,
    name VARCHAR(255)
    -- total_tickets INT,
    -- tickets_sold INT DEFAULT 0
);

CREATE TABLE tickets (
    ticket_id SERIAL PRIMARY KEY,
    event_id INT,
    status VARCHAR(20) DEFAULT 'available', -- e.g., 'available', 'sold', 'reserved'
    buyer_name VARCHAR(255) NULL,
    FOREIGN KEY (event_id) REFERENCES events(event_id)
);

INSERT INTO events (event_id, name) VALUES (1, 'Foo');

INSERT INTO tickets (event_id, status) 
SELECT 1, 'available'
FROM generate_series(1, 30); -- 30 tickets

BEGIN;
-- Select an Available Ticket:
SELECT ticket_id FROM tickets
WHERE event_id = 1 AND status = 'available'
LIMIT 1
FOR UPDATE;
-- Mark the Ticket as Sold:
UPDATE tickets 
SET status = 'sold', buyer_name = 'John Doe'
WHERE ticket_id = :ticket_id; -- ticket_id from previous SELECT (application code)
COMMIT:
```
- 2 users trying to book the same ticket, there's one would be blocked, what happen to him then after lock is released?
    - He would get an error because the ticket is already sold
        - Resolution: retry

- Retry
```sql
BEGIN;
-- Select an Available Ticket:
SELECT ticket_id FROM tickets
WHERE event_id = 1 AND status = 'available'
LIMIT 1
FOR UPDATE SKIP LOCKED;
-- Mark the Ticket as Sold:
UPDATE tickets
SET status = 'sold', buyer_name = 'John Doe'
WHERE ticket_id = :ticket_id; -- ticket_id from previous SELECT (application code)
COMMIT:
```