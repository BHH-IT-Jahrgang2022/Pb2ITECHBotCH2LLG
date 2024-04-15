package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

type Input struct {
	Text      string `json:"text"`
	SessionID string `json:"token"`
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// MOCK FUNCTION PLEASE DELETE
func mockGetSessionToken(serverURL string) string {
	// Return a predefined session token
	return "mockSessionToken"
}

//REAL GET FUNCTION UNCOMMENT!!
/*func GetSessionToken(serverURL string) string {
	// Make a GET request to /connect
	resp, err := http.Get(serverURL + "/connect")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// The body contains the session token
	sessionToken := string(body)

	return sessionToken
}*/

func sendRequest(text string, serverURL string, sessionToken string) {
	// Store the input in a struct
	input := Input{Text: text, SessionID: sessionToken}

	// Convert the struct to JSON
	jsonInput, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print debug to check if the output is correct
	fmt.Println("JSON to be sent:", string(jsonInput))

	// Create a new request using http
	req, err := http.NewRequest("GET", serverURL, bytes.NewBuffer(jsonInput))
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
// Print body of the respsonse (res.body), unmarshal the json and print json.query & json.answer
	defer res.Body.Close()
	// Print debug res-code
	if res.StatusCode != 200 {
		fmt.Println("Expected status code 200, got ", res.StatusCode)
	} else {
		fmt.Println("API call successful!")
	}
	fmt.Println(req)
}

func Initializer() {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()

	/*// Get the session token
	sessionToken := GetSessionToken(server.URL)*/

	// Mock of the above please delete afterwards
	sessionToken := mockGetSessionToken(server.URL) // Use the mock function here

	// loop this and end this with keyword
	reader := bufio.NewReader(os.Stdin)
	exitKeywords := []string{"Ende", "Exit", "Quit"} // Add more keywords as needed

	firstIteration := true
	for {
		if firstIteration {
			fmt.Println("Hello, how can I help you today?")
			firstIteration = false
		}
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if contains(exitKeywords, text) {
			fmt.Println("Ending the program.")
			break
		}

		sendRequest(text, server.URL, sessionToken) // Call the function to send the request
	}
}
