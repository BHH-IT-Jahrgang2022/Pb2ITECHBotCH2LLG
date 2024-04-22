package main

import (
	"fmt"
	"matcher/matcher"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	pingurl := "http://" + os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT") + "/data"

	timeout := 1 * time.Second

	for {
		_, err := net.DialTimeout("tcp", pingurl, timeout)

		if err != nil {
			fmt.Println("API not reachable: ", err)
		} else {
			break
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

		decoded := url.QueryEscape(req)

		fmt.Println("decoded: ", decoded)

		fmt.Println("Received request: ", req)
		c.JSON(200, gin.H{
			"response": matcher.Match(req, matches),
		})
	})
	router.GET("/reload", func(c *gin.Context) {
		matches = matcher.LoadTable()
		c.JSON(200, gin.H{
			"message": "reloaded",
		})
	})

	router.Run("0.0.0.0:8081")

}
