package main

// bombard redis til it crashes
import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// read from file file.txt
	content, err := os.ReadFile("file.txt")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// bombard redis until it crashes
	i := 1
	for {
		err := rdb.Set(ctx, fmt.Sprintf("key%d", i), string(content), 0).Err()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(i)
		i++
	}
}
