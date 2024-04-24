package main

import (
	"encoding/json"
	"fmt"
	"matcher/matcher"
	"net/http"
	"net/url"
	"os"
	"time"
	"bytes"

	"github.com/gin-gonic/gin"
)

type Unresolved struct {
	Query    string `json:"query"`
	Response string `json:"response"`
}

func main() {

	pingurl := "http://" + os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT") + "/data"

	for {
		resp, err := http.Get(pingurl)

		if err != nil {
			fmt.Println("API not reachable: ", err)
		} else {
			if resp.Body != nil {
				resp.Body.Close()
				break
			}
		}

		time.Sleep(30 * time.Second)

	}

	matches := matcher.LoadTable()

	/*
		fmt.Println("Starting Tests now...")

		test1 := "testy hat heute ??? leider blup"
		test2 := "bla mal sehen, ohne Treffer"
		test3 := "blup"
		fmt.Println(test1)
		fmt.Println(matcher.Match(test1, matches))
		fmt.Println(test2)
		fmt.Println(matcher.Match(test2, matches))
		fmt.Println(test3)
		fmt.Println(matcher.Match(test3, matches))
	*/

	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": matcher.Test(),
		})
	})
	router.GET("/match", func(c *gin.Context) {

		req := c.Query("input")
		query := c.Query("query")

		decoded_req := url.QueryEscape(req)
		decoded_query := url.QueryEscape(query)

		fmt.Println("decoded Request: ", decoded_req)
		fmt.Println("decoded Query: ", decoded_query)

		fmt.Println("Received request: ", req)
		fmt.Println("Received query: ", query)
		matched_response, resolved := matcher.Match(req, matches)
		c.JSON(200, gin.H{
			"response": matched_response,
		})
		if !resolved {
			Unresolved := Unresolved{Query: decoded_query, Response: matched_response}
			json_data, _ := json.Marshal(Unresolved)
			http.Post(os.Getenv("MAILHOST")+":"+os.Getenv("MAILPORT")+"/ticket/matchfailed", "application/json", bytes.NewBuffer(json_data))
		}
	})
	router.GET("/reload", func(c *gin.Context) {
		matches = matcher.LoadTable()
		c.JSON(200, gin.H{
			"message": "reloaded",
		})
	})

	router.Run("0.0.0.0:8081")

}
