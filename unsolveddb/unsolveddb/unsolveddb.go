package unsolveddb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

type Ticket struct {
	Tags    []string `json:"tags"`
	Problem string   `json:"problem"`
}

func FetchTicket() *[]Ticket {
	var data []Ticket
	url := "mongodb://" + os.Getenv("MONGOUSER") + ":" + os.Getenv("MONGOPASS") + "@" + os.Getenv("MONGOHOST") + ":" + os.Getenv("MONGOPORT") + "/?authSource=" + os.Getenv("MONGODB")
	fmt.Println(url)
	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if cancel == nil {
		log.Fatal(cancel)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return &[]Ticket{{[]string{"resolved"}, "test-entry"}}
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return &[]Ticket{{[]string{"resolved"}, "test-entry"}}
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return &[]Ticket{{[]string{"resolved"}, "test-entry"}}
	}

	fmt.Println("Connected to MongoDB!")

	coll := client.Database(os.Getenv("MONGODB")).Collection("tickets")

	filter := bson.D{{}}

	cur, err := coll.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return &[]Ticket{{[]string{"resolved"}, "test-entry"}}
	}

	fmt.Println("Tickets:")
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
			return &[]Ticket{}
		}

		fmt.Println(result["tags"])
		fmt.Println(result["problem"])

		strs := make([]string, len(result["tags"].(primitive.A)))

		for i, v := range result["tags"].(primitive.A) {
			strs[i] = v.(string)
		}

		data = append(data, Ticket{strs, result["problem"].(string)})

	}

	cur.Close(ctx)

	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return &[]Ticket{{[]string{"resolved"}, "test-entry"}}
	}

	fmt.Println("Connection to MongoDB closed.")

	return &data
}

func InsertTicket(data Ticket) {
	url := "mongodb://" + os.Getenv("MONGOUSER") + ":" + os.Getenv("MONGOPASS") + "@" + os.Getenv("MONGOHOST") + ":" + os.Getenv("MONGOPORT") + "/?authSource=" + os.Getenv("MONGODB")

	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if cancel == nil {
		log.Fatal(cancel)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return
	}

	fmt.Println("Connected to MongoDB!")

	collectionName := client.Database(os.Getenv("MONGODB")).Collection("tickets")

	doc := bson.D{
		{"tags", data.Tags},
		{"problem", data.Problem},
	}

	insertResult, err := collectionName.InsertOne(ctx, doc)

	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return
	}

	client.Disconnect(ctx)

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}

func UpdateTicket(data Ticket, newTags []string) string {
	url := "mongodb://" + os.Getenv("MONGOUSER") + ":" + os.Getenv("MONGOPASS") + "@" + os.Getenv("MONGOHOST") + ":" + os.Getenv("MONGOPORT") + "/?authSource=" + os.Getenv("MONGODB")
	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if cancel == nil {
		log.Fatal(cancel)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return "0"
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return "0"
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return "0"
	}

	fmt.Println("Connected to MongoDB!")

	coll := client.Database(os.Getenv("MONGODB")).Collection("tickets")

	filter := bson.D{{"problem", data.Problem}, {"tags", data.Tags}}

	cur, err := coll.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
		logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
		return "0"
	}

	var ids []string

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
			logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
			return "0"
		}

		fmt.Println(result["_id"])
		fmt.Println(result["problem"])
		fmt.Println(result["tags"])

		strs := make([]string, len(result["tags"].(primitive.A)))

		for i, v := range result["tags"].(primitive.A) {
			strs[i] = v.(string)
		}
		tempID := result["_id"].(primitive.ObjectID).Hex()
		ids = append(ids, tempID)

	}

	for _, test := range ids {
		fmt.Println(test)
		objID, _ := primitive.ObjectIDFromHex(test)
		filter := bson.D{{"_id", objID}}
		update := bson.M{
			"$set": bson.M{
				"tags": newTags,
			},
		}

		result, err := coll.UpdateOne(ctx, filter, update)

		if err != nil {
			log.Fatal(err)
			logger(LogEntry{Timestamp: time.Now().Unix(), Level: "ERROR", Message: "Error connecting to MongoDB", Service: "unsolveddb"})
			return "0"
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)

	}

	client.Disconnect(ctx)

	fmt.Println("Altered document count: ", len(ids))
	return strconv.Itoa(len(ids))
}
