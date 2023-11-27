// Using redis
// Use the SET command to cache a piece of data, e.g., SET user:1001 '{"name":"John", "age":30}'.
// Retrieve the cached data using the GET command, e.g., GET user:1001.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	user, err := json.Marshal(map[string]interface{}{
		"name":    "John",
		"age":     30,
		"married": true,
	})

	if err != nil {
		panic(err)
	}

	// 0 means the key has no expiration time.
	err = rdb.Set(ctx, "user:1", string(user), 3*time.Second).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "user:1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("user:1 string: ", val)

	type User struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}

	userObj := &User{}
	json.Unmarshal([]byte(val), userObj)
	fmt.Printf("user:1 object: %+v\n", userObj)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	time.Sleep(4 * time.Second)

	val3, err := rdb.Get(ctx, "user:1").Result()
	if err == redis.Nil {
		fmt.Println("user:1 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("user:1", val3)
	}
}
