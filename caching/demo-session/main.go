package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6380", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate user credentials

	userID := "user123" // Example user ID

	// Store session data in Redis
	err := redisClient.Set(ctx, "session_"+userID, "session_data", 0).Err()

	// Set session cookie in response:
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: userID,
	})

	if err != nil {
		http.Error(w, "Error while setting session", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Logged in user: %s", userID)
}

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get user ID from request (e.g. from cookie)
		userID := "user123" // Example user ID
		sessionData, err := redisClient.Get(ctx, "session_"+userID).Result()

		if err == redis.Nil {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Error while getting session", http.StatusInternalServerError)
			return
		}

		fmt.Println("Session Data:", sessionData)
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/login", loginHandler)
	protectedRoute := http.NewServeMux()
	protectedRoute.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Accessing protected route")
	})

	http.Handle("/", sessionMiddleware(protectedRoute))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
