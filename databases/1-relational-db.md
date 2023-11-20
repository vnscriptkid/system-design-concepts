## Exercies: https://chat.openai.com/c/a7be804f-65ff-4fd8-80cb-467e9dc9a86f
#### Atomicity with Bank Transactions:
- Exercise: Implement a bank account transfer where if the debit from one account or credit to another fails, 
neither transaction should occur.
- Guidance: Create two tables, `accounts` with `account_id`, `user_id`, `balance` (with `CHECK` constraint for non-negative balances), and `transactions` with `transaction_id`, `account_id`, `amount`, `type`, and timestamp. Use `BEGIN`, `COMMIT`, and `ROLLBACK` within a transaction block to manage the transfer.
- Answer:
```sql
CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    balance INTEGER NOT NULL CHECK (balance >= 0)
);

CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    type VARCHAR(10) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO accounts (user_id) VALUES (1);
-- ERROR:  null value in column "balance" of relation "accounts" violates not-null constraint
-- DETAIL:  Failing row contains (1, 1, null).
INSERT INTO accounts (user_id, balance) VALUES (1, -1);
-- ERROR:  new row for relation "accounts" violates check constraint "accounts_balance_check"
-- DETAIL:  Failing row contains (2, 1, -1).
INSERT INTO accounts (account_id, user_id, balance) VALUES (1, 1, 0);

-- Top up account 1 with $100
BEGIN;
    INSERT INTO transactions ("account_id", "amount", "type") VALUES (1, 100, 'credit');
    UPDATE "accounts" SET "balance" = "balance" + 100 WHERE account_id = 1;
COMMIT;
-- Withdraw $150 from account 1
BEGIN;
    INSERT INTO transactions ("account_id", "amount", "type") VALUES (1, 150, 'debit');
    UPDATE "accounts" SET "balance" = "balance" - 150 WHERE account_id = 1;
    -- ERROR:  new row for relation "accounts" violates check constraint "accounts_balance_check"
    -- DETAIL:  Failing row contains (1, 1, -50).
    SELECT * FROM transactions;
    -- ERROR:  current transaction is aborted, commands ignored until end of transaction block
COMMIT;
-- ROLLBACK

```

#### Consistency with Account Creation through CHECK Constraint:
- Exercise: Enforce data integrity when creating user accounts with complete and valid information.
- Guidance: Define a `users` table with `user_id`, `username` (UNIQUE), `email` (UNIQUE and with a REGEX pattern constraint for format validation), password (with a CHECK constraint for complexity), and `created_at` (DEFAULT to the current timestamp). Use INSERT statements with constraints that enforce consistency.
- Answer:
```sql
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(20) NOT NULL
        CONSTRAINT users_username_unique UNIQUE,
    email VARCHAR(255) NOT NULL
        CONSTRAINT users_email_unique UNIQUE
        CONSTRAINT users_email_regex CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    "password" VARCHAR(255) NOT NULL
        -- At least 8 characters, 1 uppercase, 1 lowercase, 1 number
        CONSTRAINT users_password_complexity CHECK ("password" ~* '^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$'),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, email, "password") VALUES ('user1', 'user1@gmail.com', 'Password1');
-- INSERT 0 1
-- Fail uniqueness constraint on username
INSERT INTO users (username, email, "password") VALUES ('user1', 'user1plus@gmail.com', 'Password1');
-- ERROR:  duplicate key value violates unique constraint "users_username_unique"
-- DETAIL:  Key (username)=(user1) already exists.

-- Fail email format constraint
INSERT INTO users (username, email, "password") VALUES ('user2', 'user2gmail.com', 'Password1');
-- ERROR:  new row for relation "users" violates check constraint "users_email_regex"
-- DETAIL:  Failing row contains (3, user2, user2gmail.com, Password1, 2023-11-19 17:07:28.022551).

-- Fail primary key constraint
INSERT INTO users (user_id, username, email, "password") VALUES (100, 'user100', 'user100@gmail.com', 'Password1');
INSERT INTO users (user_id, username, email, "password") VALUES (100, 'user101', 'user101@gmail.com', 'Password1');
-- ERROR:  duplicate key value violates unique constraint "users_pkey"
-- DETAIL:  Key (user_id)=(100) already exists.

-- Implicit transaction in batch insert
INSERT INTO users (username, email, "password") VALUES 
    ('user200', 'user200@gmail.com', 'Password1'), -- valid
    ('user201', 'user201@gmail.c', 'Password1'); -- invalid
```

#### Isolation with Concurrent Access:
- Exercise: Two customers attempt to book the last ticket for a flight, but the system should prevent double booking.
- Guidance: Use a flights table with flight_id, seats_available (with a CHECK constraint to prevent negative numbers). Utilize SERIALIZABLE isolation level to prevent phantom reads and ensure that once one customer books the seat, the other cannot.

#### Durability with Order Processing:
- Exercise: Guarantee that once an order is processed, it's permanently stored even if a crash occurs immediately after.
- Guidance: Create an orders table with order_id, user_id, order_details, and status. After inserting the order, simulate a failure (e.g., stop the database service). Restart and query the orders table to ensure the data persists due to the durability property.

#### Atomicity with Batch Processing:
- Exercise: Process a batch of data insertions where one insertion failure should prevent all insertions.
- Guidance: Create a batch_process table and use a transaction block to insert multiple records. Utilize TRY/CATCH blocks for error handling and use ROLLBACK on error detection to maintain atomicity.
```sql
CREATE TABLE batch_process (
    id SERIAL PRIMARY KEY,
    data VARCHAR(255) NOT NULL
);
```

```go
func batchInsert() error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    for _, data := range dataSlice {
        _, err := tx.Exec("INSERT INTO batch_process(data) VALUES (?)", data)
        if err != nil {
            tx.Rollback()
            return err
        }
    }

    tx.Commit()
    return nil
}
```

#### Consistency with Data Update through CHECK Constraint:
- Exercise: Update user profiles ensuring all emails follow a specific format.
- Guidance: In the users table, apply a CHECK constraint on the email field to enforce a specific format using REGEX. Test updates with a mix of valid and invalid emails to see consistency in action.

```sql
INSERT INTO users (username, email, "password") VALUES ('user300', 'user300@gmail.com', 'Password1');
-- INSERT 0 1
UPDATE users SET email = 'user300gmail.com' WHERE username = 'user300';
-- ERROR:  new row for relation "users" violates check constraint "users_email_regex"
-- DETAIL:  Failing row contains (8, user300, user300gmail.com, Password1, 2023-11-19 17:53:48.92215).
UPDATE users SET email = 'user300plus@gmail.com' WHERE username = 'user300';
-- UPDATE 1
```

#### Isolation with Simulated Deadlock:
- Exercise: Simulate a deadlock scenario where two processes lock each other out.
- Guidance: Create a deadlock_test table with resource_id and lock_status. Begin two transactions that attempt to update the same row. Use SQL tools to monitor and resolve the deadlock. Understand how isolation levels can prevent this.

#### Durability with Transaction Log:
- Exercise: Perform an insert operation and then analyze the transaction log.
- Guidance: Insert data into a log_test table. Then, query the transaction log (the specific query will depend on your DBMS) to find the log entry for your insert, showing how durability is logged.

#### Atomicity with Partial Success:
- Exercise: Run a multi-step data entry process where one step fails intentionally.
- Guidance: In a multi_step table with appropriate data types and constraints, use a transaction to group insertions. Include an erroneous insertion to trigger a failure. Observe the atomic rollback.

#### Consistency with Referential Integrity:
- Exercise: Add data to a table with a foreign key constraint to ensure referential integrity.
- Guidance: Create products and orders tables where orders includes a product_id foreign key that references product_id in products. Attempt to insert into orders with both valid and invalid product_id values to test the foreign key constraint.
```sql
CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(product_id)
    -- ON DELETE RESTRICT (default) >>> strict
    -- ON DELETE CASCADE >>> so brutal
    -- ON DELETE SET NULL >>> orphan records (no parent), won't work with `NOT NULL`
);

INSERT INTO products (product_id, name) VALUES (1, 'Product 1');
-- INSERT 0 1
INSERT INTO orders (product_id) VALUES (1);
-- INSERT 0 1
INSERT INTO orders (product_id) VALUES (2);
-- ERROR:  insert or update on table "orders" violates foreign key constraint "orders_product_id_fkey"
-- DETAIL:  Key (product_id)=(2) is not present in table "products".

-- Delete product 1
DELETE FROM products WHERE product_id = 1;
-- ERROR:  update or delete on table "products" violates foreign key constraint "orders_product_id_fkey" on table "orders"
-- DETAIL:  Key (product_id)=(1) is still referenced from table "orders"

DROP TABLE products;
-- ERROR:  cannot drop table products because other objects depend on it
-- DETAIL:  constraint orders_product_id_fkey on table orders depends on table products
-- HINT:  Use DROP ... CASCADE to drop the dependent objects too.
```

#### Atomicity in Refunding Transactions:
- Exercise: Simulate a refund process where if any part of the transaction fails, no changes should take effect.
- Guidance: Use a sales table with sale_id, customer_id, amount, and status. Also, have a refunds table with refund_id, sale_id, and amount. Wrap the refund update and insertion into refunds table within a transaction. If the insertion fails, ensure the update is not committed.

#### Consistency in Library Management:
- Exercise: Add a new book entry ensuring that all bibliographic fields meet specific standards.
- Guidance: Define a books table with isbn (UNIQUE), title, author, publication_date, pages (integer, CHECK > 0), and available_copies (integer, CHECK >= 0). Insert new books using constraints to maintain data consistency.

#### Isolation in Auction Bidding:
- Exercise: Manage simultaneous bids on an item so that they do not interfere with each other.
- Guidance: Create an auctions table with item_id, current_bid, and bidder_id. Use transactions with the REPEATABLE READ isolation level when updating bids to prevent lost updates or dirty reads.

#### Durability in Customer Reviews:
- Exercise: Ensure that customer reviews are not lost even if the system crashes right after submission.
- Guidance: Create a reviews table with review_id, customer_id, product_id, rating, comment, and timestamp. After inserting a review, simulate a system crash. Upon recovery, verify the presence of the review to demonstrate durability.

#### Atomicity in Inventory Adjustment:
- Exercise: Adjust inventory levels for multiple products, ensuring all adjustments are applied or none at all.
- Guidance: Define an inventory table with product_id, quantity_available (integer, CHECK >= 0). Perform a series of updates within a transaction, and simulate an error on one update to test atomicity with a rollback.

#### Consistency in Course Enrollment:
- Exercise: Enforce rules when students enroll in courses, such as prerequisites and maximum class size.
- Guidance: Create courses and enrollments tables with course_id, student_id, enrolled_on date, and status. Add CHECK constraints on status for valid values, and use a trigger to check for prerequisite courses before allowing enrollment.
```sql
CREATE TABLE courses (
    course_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    prerequisite_course_id INTEGER REFERENCES courses(course_id)
);

CREATE TABLE enrollments (
    enrollment_id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(course_id),
    student_id INTEGER NOT NULL,
    enrolled_on DATE NOT NULL,
    status VARCHAR(10) NOT NULL
        CONSTRAINT enrollments_status_check CHECK (status IN ('enrolled', 'dropped'))
);

-- setup trigger
CREATE OR REPLACE FUNCTION check_prerequisite() RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM courses
        WHERE course_id = NEW.course_id
        -- This new course has a prerequisite and the student has not taken it
        AND prerequisite_course_id IS NOT NULL
        AND prerequisite_course_id NOT IN (
            SELECT course_id FROM enrollments
            WHERE student_id = NEW.student_id
            AND status = 'enrolled'
        )
    ) THEN
        RAISE EXCEPTION 'Prerequisite not met';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_prerequisite
BEFORE INSERT ON enrollments -- Even using AFTER INSERT, the error is still thrown and the row is not inserted
FOR EACH ROW
EXECUTE PROCEDURE check_prerequisite();

-- insert courses
INSERT INTO courses (course_id, name) VALUES (1, 'Course 1');
INSERT INTO courses (course_id, name) VALUES (2, 'Course 2');
INSERT INTO courses (course_id, name, prerequisite_course_id) VALUES (3, 'Course 3', 1);
INSERT INTO courses (course_id, name, prerequisite_course_id) VALUES (4, 'Course 4', 2);

-- insert enrollments
INSERT INTO enrollments (course_id, student_id, enrolled_on, status) VALUES (1, 1, '2021-01-01', 'enrolled');
INSERT INTO enrollments (course_id, student_id, enrolled_on, status) VALUES (3, 1, '2021-01-01', 'enrolled');
INSERT INTO enrollments (course_id, student_id, enrolled_on, status) VALUES (4, 1, '2021-01-01', 'enrolled');
-- ERROR:  Prerequisite not met

```

#### Isolation in Financial Ledger Entries:
- Exercise: Enter multiple financial transactions in a ledger where each entry must be isolated to maintain accuracy.
- Guidance: In a financial_ledger table with transaction_id, account_id, debit_amount, credit_amount, perform simultaneous debit and credit operations using the SERIALIZABLE isolation level to ensure accurate balance sheets.

#### Durability in User Session Logs:
- Exercise: Record user login and logout times, ensuring that session logs are always recorded.
- Guidance: Use a session_logs table with log_id, user_id, login_time, logout_time. After logging the events, terminate the database process and then restart to check log integrity.

#### Atomicity in Batch User Registrations:
- Exercise: Process a batch of user registrations where failure in any registration should nullify the entire batch.
- Guidance: Create a batch_registrations table and wrap multiple INSERT statements in a transaction. Use exception handling to catch any insertion errors and roll back if necessary.

#### Consistency in Hotel Bookings:
- Exercise: Ensure all hotel bookings adhere to business rules, such as room availability and booking dates.
- Guidance: Set up a hotel_rooms table with room_number, room_type, is_booked (boolean), and booking_dates (date range). Use constraints to ensure is_booked accurately reflects the booking_dates, and no overlapping bookings can occur.
```sql
CREATE TABLE hotel_rooms (
    room_number INTEGER PRIMARY KEY,
    room_type VARCHAR(255) NOT NULL,
    is_booked BOOLEAN NOT NULL DEFAULT FALSE,
    booking_dates DATERANGE NOT NULL
        CONSTRAINT hotel_rooms_booking_dates_check CHECK (NOT is_booked OR booking_dates IS NOT NULL)
);

INSERT INTO hotel_rooms (room_number, room_type, is_booked, booking_dates) VALUES (1, 'Single', TRUE, '[2021-01-01, 2021-01-02)');
INSERT INTO hotel_rooms (room_number, room_type, is_booked, booking_dates) VALUES (2, 'Single', false, '[2021-01-01, 2021-01-02)');

```

#### Trigger for Automatic Balance Updates
- Exercise: Create a trigger on a banking system's transactions table that automatically updates the account balance whenever a new transaction is recorded.
- Guidance: Define a transactions table with columns like transaction_id, account_id, amount, and type (debit or credit). Then, create a trigger that fires after an insert operation on this table. The trigger should update the balance column in the accounts table, adding the amount for credit transactions and subtracting for debit transactions, ensuring the balance is always accurate.

```sql
CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    balance INTEGER NOT NULL CHECK (balance >= 0)
);

CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(account_id),
    amount INTEGER NOT NULL CHECK (amount > 0),
    type VARCHAR(10) NOT NULL
);

CREATE OR REPLACE FUNCTION update_balance() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.type = 'credit' THEN
        UPDATE accounts SET balance = balance + NEW.amount WHERE account_id = NEW.account_id;
    ELSIF NEW.type = 'debit' THEN
        UPDATE accounts SET balance = balance - NEW.amount WHERE account_id = NEW.account_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_balance
AFTER INSERT ON transactions
FOR EACH ROW
EXECUTE PROCEDURE update_balance();

INSERT INTO accounts (account_id, balance) VALUES (1, 0);
INSERT INTO transactions (account_id, amount, type) VALUES (1, 100, 'credit');
INSERT INTO transactions (account_id, amount, type) VALUES (1, 50, 'debit');
INSERT INTO transactions (account_id, amount, type) VALUES (1, 55, 'debit');
-- ERROR:  new row for relation "accounts" violates check constraint "accounts_balance_check"
-- DETAIL:  Failing row contains (1, -5).
-- CONTEXT:  SQL statement "UPDATE accounts SET balance = balance - NEW.amount WHERE account_id = NEW.account_id"
-- PL/pgSQL function update_balance() line 6 at SQL statement

```

#### Trigger for Order and Inventory Management
- Exercise: Implement a trigger to adjust inventory levels automatically when a new order is placed.
- Guidance: Use an orders table with order_id, product_id, and quantity. The inventory table should have product_id and quantity_in_stock. Create an AFTER INSERT trigger on the orders table that reduces the quantity_in_stock in the inventory table by the ordered quantity, maintaining consistent inventory records.

#### Trigger for Referential Integrity Checks
- Exercise: Design a trigger to prevent deletion of a customer record if they have existing orders.
- Guidance: Have a customers table and an orders table with a foreign key relationship. Create a BEFORE DELETE trigger on customers that checks if there are any related records in the orders table. If there are, the trigger should prevent the deletion, thereby maintaining referential integrity.

#### Trigger for Enforcing Complex Business Rules
- Exercise: Use a trigger to enforce a business rule where an employee's salary cannot exceed a certain limit based on their department.
- Guidance: Create an employees table with employee_id, department_id, and salary. Also, have a department table with department_id and salary_limit. Write a BEFORE INSERT or BEFORE UPDATE trigger on the employees table that compares the employee's salary with the salary_limit of their department, ensuring the consistency of salary distribution across the company.

#### Trigger for Synchronizing Redundant Data
- Exercise: Create a trigger to keep redundant data synchronized across different tables.
- Guidance: Suppose you have two tables, employees and department_managers, where department_managers stores a redundant copy of the manager's name. Write an AFTER UPDATE trigger on the employees table so that when a manager's name changes, it automatically updates the corresponding record in department_managers. This ensures that redundant data remains consistent across tables.
```sql
```