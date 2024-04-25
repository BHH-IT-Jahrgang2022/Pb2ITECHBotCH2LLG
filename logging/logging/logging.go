package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
}

var ctx = context.Background()

// Set a log entry in the database
func set(c *redis.Client, key int64, value LogEntry) {
	p, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Error marshalling value ", err)
	}
	c.Set(ctx, fmt.Sprintf("%d", key), p, 0)
}

// Get a log entry from the database
func get(c *redis.Client, key string) LogEntry {
	p := c.Get(ctx, key).Val()
	fmt.Println("Value from redis ", p)
	var log LogEntry
	err := json.Unmarshal([]byte(p), &log)
	if err != nil {
		fmt.Println("Error unmarshalling value ", err)
	}
	return log
}

// Initialize the Redis client
func InitClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
	return rdb
}

func StartApi() {
	db := InitClient()
	if db == nil {
		db_connection_timeout := 12
		for i := 0; i < db_connection_timeout; i++ {
			time.Sleep(10 * time.Second)
			db = InitClient()
			if db != nil {
				break
			}
		}
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// Endpoint for services to send logs to
	r.POST("/log", func(c *gin.Context) {
		var log LogEntry
		if err := c.ShouldBindJSON(&log); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		set(db, time.Now().UnixNano(), log)
		c.JSON(http.StatusOK, gin.H{"message": "Log saved"})
	})
	// Returns the entire log, only for debugging purposes
	r.GET("/dump", func(c *gin.Context) {
		logs := make(map[string]LogEntry)
		keys, _ := db.Keys(ctx, "*").Result()
		for _, key := range keys {
			entry := get(db, key)
			logs[key] = entry
		}
		fmt.Println("Logs ", logs)
		c.JSON(http.StatusOK, gin.H{"log_length": len(logs)})
	})
	r.GET("")
	r.Run("0.0.0.0:8080")
}
