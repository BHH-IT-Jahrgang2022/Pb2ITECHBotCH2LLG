package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test() string {
	return "I'm alive"
}

type Data struct {
	Keywords []string `json:"keywords"`
	Response string   `json:"response"`
}

func FetchData() *[]Data {
	var data []Data
	clientOptions := options.Client().ApplyURI("mongodb://user:password@127.0.0.1:27017/?authSource=db")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if cancel == nil {
		log.Fatal(cancel)
		return &[]Data{}
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return &[]Data{}
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return &[]Data{}
	}

	fmt.Println("Connected to MongoDB!")

	collectionWindowfly := client.Database("test").Collection("windowfly")
	collectionCleanbug := client.Database("test").Collection("cleanbug")
	collectionGardenbeetle := client.Database("test").Collection("gardenbeetle")
	collectionEmpty := client.Database("test").Collection("empty")

	filter := bson.D{{}}

	curFly, err := collectionWindowfly.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
		return &[]Data{}
	}

	curBug, err := collectionCleanbug.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return &[]Data{}
	}
	curBeetle, err := collectionGardenbeetle.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return &[]Data{}
	}
	curEmpty, err := collectionEmpty.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return &[]Data{}
	}

	fmt.Println("Windowfly:")
	for curFly.Next(ctx) {
		var result bson.M
		err := curFly.Decode(&result)
		if err != nil {
			log.Fatal(err)
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
		return &[]Data{}
	}

	fmt.Println("Connection to MongoDB closed.")

	return &data
}
