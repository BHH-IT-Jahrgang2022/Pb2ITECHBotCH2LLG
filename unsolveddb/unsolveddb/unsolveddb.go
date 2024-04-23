package unsolveddb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if cancel == nil {
		log.Fatal(cancel)
		return &[]Ticket{{[]string{"resolved"}, "test-entry"}}
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return &[]Ticket{{[]string{"resolved"}, "test-entry"}}
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return &[]Ticket{{[]string{"resolved"}, "test-entry"}}
	}

	fmt.Println("Connected to MongoDB!")

	coll := client.Database(os.Getenv("MONGODB")).Collection("tickets")

	filter := bson.D{{}}

	cur, err := coll.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
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
		return
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
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
		return
	}

	client.Disconnect(ctx)

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}
