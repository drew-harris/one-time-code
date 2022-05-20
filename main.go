package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var rdb *redis.Client

func handle(w http.ResponseWriter, req *http.Request) {
	code := strings.TrimPrefix(req.URL.Path, "/code/")
	fmt.Println(code)

	_, err := rdb.Get(context.Background(), code).Result()

	if err != nil {
		rdb.Set(context.Background(), code, "VALUE", time.Hour*999*999).Result()
		w.Write([]byte("Success"))
		return
	}

	w.WriteHeader(404)
	w.Write([]byte("Error"))
}

func main() {
	// Connect to redis
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not get environment variables")
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Get a test value to ensure connection
	ping, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(ping)

	// Set up HTTP server
	fmt.Println("Starting HTTP Server")
	http.HandleFunc("/code/", handle)
	http.ListenAndServe(":8080", nil)
}
