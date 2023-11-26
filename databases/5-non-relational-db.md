## Cover
- Document DB
    - MongoDB
    - ElasticSearch
- Key-value DB
    - Redis
    - DynamoDB

## Doubts
- Partial update in MongoDB
    - Example: update only the `name` field of a document
    - Solution: use `$set` operator
    ```js
    db.users.update(
        { _id: 1 },
        { $set: { name: "John Doe" } }
    )
    ```
    - Compared to relational DB:
        - `UPDATE users SET name = "John Doe" WHERE id = 1;`
        - Difference:
            - relational DB: update the whole record
            - MongoDB: update only the `name` field
- Document DB use cases:
    - in-app notifications
        - Why?
            - unstructured data
        - Example full-fledged json record:
        ```json
        {
            "user_id": "123",
            "title": "New message",
            "body": "You have a new message from John Doe",
            "timestamp": "2020-01-01 12:00:00",
            "read": false,
        }
        ```

    - catalog service
        - Why?
            - unstructured data
        - Example full-fledged json record:
        ```json
        {
            "product_id": "123",
            "name": "iPhone 12",
            "description": "Apple's latest iPhone",
            "price": 1000,
            "category": "Electronics",
            "tags": ["smartphone", "apple", "iphone"],
            "reviews": [
                {
                    "user_id": "123",
                    "rating": 5,
                    "comment": "Great phone!"
                },
                {
                    "user_id": "456",
                    "rating": 4,
                    "comment": "Good phone!"
                }
            ]
        }
        ```

- Key-value DB?
    - When to use?
        - fast read/write
        - access data by key
    - Use-cases:
        - order management
            - access order by order_id
            - what if i want get all orders by user_id?
                - use secondary index (user_id): https://youtu.be/ihMOlb8EZKE?si=wIWVAjHk0zqWYfj0
        - user session management
        - shopping cart
        - caching
        - leaderboard
