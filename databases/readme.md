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



