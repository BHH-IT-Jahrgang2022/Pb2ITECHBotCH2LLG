package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
}

func logger(entry LogEntry) {
	// Send the log entry to the logging API
	if os.Getenv("LOGGING_ENABLED") == "true" {
		logging_API_route := os.Getenv("LOGGING_API_ROUTE")
		jsonEntry, _ := json.Marshal(entry)
		http.Post(logging_API_route+"/log", "application/json", bytes.NewBuffer(jsonEntry))
	}
}

func Test() string {
	return "I'm alive"
}

type Data struct {
	Keywords []string `json:"keywords"`
	Response string   `json:"response"`
}

func FetchData() *[]Data {
	var data []Data
	url := "mongodb://" + os.Getenv("MONGOUSER") + ":" + os.Getenv("MONGOPASS") + "@" + os.Getenv("MONGOHOST") + ":" + os.Getenv("MONGOPORT") + "/?authSource=" + os.Getenv("MONGODB")
	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if cancel == nil {
		log.Fatal(cancel)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in context.WithTimeout", Service: "db"})
		return &[]Data{}
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in mongo.Connect", Service: "db"})
		return &[]Data{}
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in client.Ping", Service: "db"})
		return &[]Data{}
	}

	fmt.Println("Connected to MongoDB!")

	collectionWindowfly := client.Database(os.Getenv("MONGODB")).Collection("windowfly")
	collectionCleanbug := client.Database(os.Getenv("MONGODB")).Collection("cleanbug")
	collectionGardenbeetle := client.Database(os.Getenv("MONGODB")).Collection("gardenbeetle")
	collectionEmpty := client.Database(os.Getenv("MONGODB")).Collection("empty")

	filter := bson.D{{}}

	curFly, err := collectionWindowfly.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in collectionWindowfly.Find", Service: "db"})
		return &[]Data{}
	}

	curBug, err := collectionCleanbug.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in collectionCleanbug.Find", Service: "db"})
		return &[]Data{}
	}
	curBeetle, err := collectionGardenbeetle.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in collectionGardenbeetle.Find", Service: "db"})
		return &[]Data{}
	}
	curEmpty, err := collectionEmpty.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in collectionEmpty.Find", Service: "db"})
		return &[]Data{}
	}

	fmt.Println("Windowfly:")
	for curFly.Next(ctx) {
		var result bson.M
		err := curFly.Decode(&result)
		if err != nil {
			log.Fatal(err)
			logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in curFly.Decode", Service: "db"})
			return &[]Data{}
		}

		fmt.Println(result["keywords"])
		fmt.Println(result["response"])

		strs := make([]string, len(result["keywords"].(primitive.A)))

		for i, v := range result["keywords"].(primitive.A) {
			strs[i] = v.(string)
		}

		data = append(data, Data{strs, result["response"].(string)})

	}

	fmt.Println("Cleanbug:")
	for curBug.Next(ctx) {
		var result bson.M
		err := curBug.Decode(&result)
		if err != nil {
			log.Fatal(err)
			logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in curBug.Decode", Service: "db"})
			return &[]Data{}
		}
		fmt.Println(result["keywords"])
		fmt.Println(result["response"])

		strs := make([]string, len(result["keywords"].(primitive.A)))

		for i, v := range result["keywords"].(primitive.A) {
			strs[i] = v.(string)
		}

		data = append(data, Data{strs, result["response"].(string)})

	}

	fmt.Println("Gardenbeetle:")
	for curBeetle.Next(ctx) {
		var result bson.M
		err := curBeetle.Decode(&result)
		if err != nil {
			log.Fatal(err)
			logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in curBeetle.Decode", Service: "db"})
			return &[]Data{}
		}
		fmt.Println(result["keywords"])
		fmt.Println(result["response"])

		strs := make([]string, len(result["keywords"].(primitive.A)))

		for i, v := range result["keywords"].(primitive.A) {
			strs[i] = v.(string)
		}

		data = append(data, Data{strs, result["response"].(string)})
	}

	fmt.Println("Empty:")
	for curEmpty.Next(ctx) {
		var result bson.M
		err := curEmpty.Decode(&result)
		if err != nil {
			log.Fatal(err)
			logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in curEmpty.Decode", Service: "db"})
			return &[]Data{}
		}
		fmt.Println(result["keywords"])
		fmt.Println(result["response"])
		strs := make([]string, len(result["keywords"].(primitive.A)))

		for i, v := range result["keywords"].(primitive.A) {
			strs[i] = v.(string)
		}

		data = append(data, Data{strs, result["response"].(string)})
	}

	curFly.Close(ctx)
	curBug.Close(ctx)
	curBeetle.Close(ctx)
	curEmpty.Close(ctx)

	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in client.Disconnect", Service: "db"})
		return &[]Data{}
	}

	fmt.Println("Connection to MongoDB closed.")

	return &data
}

func InsertData(data Data, collection string) {
	url := "mongodb://" + os.Getenv("MONGOUSER") + ":" + os.Getenv("MONGOPASS") + "@" + os.Getenv("MONGOHOST") + ":" + os.Getenv("MONGOPORT") + "/?authSource=" + os.Getenv("MONGODB")
	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if cancel == nil {
		log.Fatal(cancel)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in context.WithTimeout", Service: "db"})
		return
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in mongo.Connect", Service: "db"})
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in client.Ping", Service: "db"})
		return
	}

	fmt.Println("Connected to MongoDB!")

	collectionName := client.Database(os.Getenv("MONGODB")).Collection(collection)

	doc := bson.D{
		{"keywords", data.Keywords},
		{"response", data.Response},
	}

	insertResult, err := collectionName.InsertOne(ctx, doc)

	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "Error", Message: "Error in collectionName.InsertOne", Service: "db"})
		return
	}

	client.Disconnect(ctx)

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}
