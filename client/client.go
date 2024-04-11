package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
)

type Input struct {
	Text string `json:"text"`
}

func main() {
	// Mock server 
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Hello, how can i help you today?:")

	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1] // Remove newline character

	// Store the input in a struct
	input := Input{Text: text}

	// Convert the struct to JSON
	jsonInput, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print debug to check if the output is correct
	fmt.Println("JSON to be sent:", string(jsonInput))

	// Create a new request using http -> Replace "server.URL" with the actual URL of the Endpoint
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonInput))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the content type header value
	req.Header.Set("Content-Type", "application/json")

	// Send the request via a client
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Expected status code 200, got ", res.StatusCode)
	} else {
		fmt.Println("API call successful!")
	}
	fmt.Println(req)

}
