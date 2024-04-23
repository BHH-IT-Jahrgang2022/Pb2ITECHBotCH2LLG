package tokenizer

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestMatcherResponse(t *testing.T) {

	matcherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"response": "hallo auch"}`)
	}))
	defer matcherServer.Close()

	os.Setenv("MATCHER_URL", matcherServer.URL)

	var matcherURL string
	oldMatcherURL := matcherURL
	matcherURL = matcherServer.URL
	defer func() { matcherURL = oldMatcherURL }()

	query := "hello world!@#$%^&*()"
	escapedQuery := url.QueryEscape(query)
	req, err := http.NewRequest("GET", "http://localhost:8080/tokenize?query="+escapedQuery, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"response": "hallo auch"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestTokenize(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"hello world", []string{"hello", "world"}},
		{"hello, world!", []string{"hello", "world"}},
		{"", []string{}},
		{"  $§& ", []string{}},
		{"hello   world", []string{"hello", "world"}},
		{"hello *§$%(§&)world!", []string{"hello", "world"}},
	}

	for _, testCase := range testCases {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in test", r)
			}
		}()

		fmt.Printf("Input: %v\n", testCase.input)
		output := Tokenize(testCase.input)
		fmt.Printf("Output: %v\n", output)

		fmt.Printf("Output   >>> Type: %T, Length: %d, Capacity: %d, Value: %v\n", output, len(output), cap(output), output)
		fmt.Printf("Expected >>> Type: %T, Length: %d, Capacity: %d, Value: %v\n", testCase.expected, len(testCase.expected), cap(testCase.expected), testCase.expected)

		if output == nil {
			t.Errorf("Output should not be nil")
		}

		if testCase.expected == nil {
			t.Errorf("Expected output should not be nil")
		} else if len(output) == 0 && len(testCase.expected) == 0 {
			// Both slices are empty, so they are equal
			continue
		} else if !reflect.DeepEqual(output, testCase.expected) {
			fmt.Printf("Actual: %v, Expected: %v\n", output, testCase.expected)
			t.Errorf("Tokenize(%q) = %v, want %v", testCase.input, output, testCase.expected)
		}
	}
}

func TestMissingResponseAttribute(t *testing.T) {
	// Create a new server to mock the matcher which returns a JSON response without the expected attribute
	matcherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"not_response": "wrong attribute"}`)
	}))
	defer matcherServer.Close()

	// Set the MATCHER_URL environment variable to the URL of the mock matcher server
	os.Setenv("MATCHER_URL", matcherServer.URL)
	var matcherURL string
	oldMatcherURL := matcherURL
	matcherURL = matcherServer.URL
	defer func() { matcherURL = oldMatcherURL }()

	// Create a new request to the tokenizer server with a query that will be tokenized
	query := "hello world!@#$%^&*()"
	escapedQuery := url.QueryEscape(query)
	req, err := http.NewRequest("GET", "http://localhost:8080/tokenize?query="+escapedQuery, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to record the response from the tokenizer server
	rr := httptest.NewRecorder()

	// Create a new handler from the handleRequest function
	handler := http.HandlerFunc(handleRequest)

	// Send the request to the tokenizer server
	handler.ServeHTTP(rr, req)

	// Check if the status code of the response is 400 (Bad Request)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check if the body of the response is the expected error message
	expected := "Missing expected attribute in JSON from matcher"
	received := strings.TrimSpace(rr.Body.String())
	fmt.Println("Body.String (trimmed): " + received)
	if received != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", received, expected)
	}
}

func TestWrongQuantityOfAttributes(t *testing.T) {
	// Create a new server to mock the matcher which returns a JSON response with more than one attribute
	matcherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"response": "hallo auch", "extra": "attribute"}`)
	}))
	defer matcherServer.Close()

	// Set the MATCHER_URL environment variable to the URL of the mock matcher server
	os.Setenv("MATCHER_URL", matcherServer.URL)
	var matcherURL string
	oldMatcherURL := matcherURL
	matcherURL = matcherServer.URL
	defer func() { matcherURL = oldMatcherURL }()

	// Create a new request to the tokenizer server with a query that will be tokenized
	query := "hello world!@#$%^&*()"
	escapedQuery := url.QueryEscape(query)
	req, err := http.NewRequest("GET", "http://localhost:8080/tokenize?query="+escapedQuery, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to record the response from the tokenizer server
	rr := httptest.NewRecorder()

	// Create a new handler from the handleRequest function
	handler := http.HandlerFunc(handleRequest)

	// Send the request to the tokenizer server
	handler.ServeHTTP(rr, req)

	// Check if the status code of the response is 400 (Bad Request)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check if the body of the response is the expected error message
	expected := "Expected 1 attributes, got 2"
	received := strings.TrimSpace(rr.Body.String())
	if received != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", received, expected)
	}
}

func TestInvalidJsonFromMatcher(t *testing.T) {
	// Create a server to mock the matcher which returns an invalid response like a string or an array
	matcherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "invalid JSON")
	}))
	defer matcherServer.Close()

	// Set the MATCHER_URL environment variable to the URL of the mock matcher server
	os.Setenv("MATCHER_URL", matcherServer.URL)
	var matcherURL string
	oldMatcherURL := matcherURL
	matcherURL = matcherServer.URL
	defer func() { matcherURL = oldMatcherURL }()

	// Create a new request to the tokenizer server with a query that will be tokenized
	query := "hello world!@#$%^&*()"
	escapedQuery := url.QueryEscape(query)
	req, err := http.NewRequest("GET", "http://localhost:8080/tokenize?query="+escapedQuery, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to record the response from the tokenizer server
	rr := httptest.NewRecorder()

	// Create a new handler from the handleRequest function
	handler := http.HandlerFunc(handleRequest)

	// Send the request to the tokenizer server
	handler.ServeHTTP(rr, req)

	// Check if the status code of the response is 500 (Internal Server Error)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check if the body of the response is the expected error message
	expected := "Received invalid JSON from matcher"
	received := strings.TrimSpace(rr.Body.String())
	if received != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", received, expected)
	}

}

func TestHttpMethod(t *testing.T) {
	tests := []struct {
		method string
		status int
	}{
		{"POST", http.StatusMethodNotAllowed},
		{"PUT", http.StatusMethodNotAllowed},
		{"DELETE", http.StatusMethodNotAllowed},
		{"PATCH", http.StatusMethodNotAllowed},
		{"GET", http.StatusOK},
	}

	// Create a server to mock the matcher which returns an invalid response like a string or an array
	matcherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"response": "hallo auch"}`)
	}))
	defer matcherServer.Close()

	// Set the MATCHER_URL environment variable to the URL of the mock matcher server
	os.Setenv("MATCHER_URL", matcherServer.URL)
	var matcherURL string
	oldMatcherURL := matcherURL
	matcherURL = matcherServer.URL
	defer func() { matcherURL = oldMatcherURL }()

	for _, test := range tests {
		req, err := http.NewRequest(test.method, "http://localhost:8080/tokenize?query=hello", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleRequest)

		handler.ServeHTTP(rr, req)

		fmt.Println(rr.Result())
		fmt.Println(rr.Body.String())

		if status := rr.Code; status != test.status {
			t.Errorf("handler returned wrong status code: got %v want %v", status, test.status)
		}
	}

}

func TestMissingQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/tokenize", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := "Missing query parameter"
	received := strings.TrimSpace(rr.Body.String())
	if received != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleRequest_ErrorStatus(t *testing.T) {
	// Erstellen Sie einen HTTP-Testserver, der einen Fehlerstatus zurückgibt
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	// Schließen Sie den Server, wenn der Test beendet ist
	defer server.Close()

	// Ersetzen Sie die matcher_URL durch die URL des Test-Servers
	matcher_URL := server.URL

	// Führen Sie die Funktion aus, die Sie testen möchten
	resp, err := http.Get(matcher_URL + "/match?input=" + "test")

	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	// Überprüfen Sie, ob der Statuscode der Antwort dem erwarteten Wert entspricht
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)

	}
}
