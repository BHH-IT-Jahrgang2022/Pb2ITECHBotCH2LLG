package tokenizer

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
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
		output := Tokenize(testCase.input)

		fmt.Printf("Output   >>> Type: %T, Length: %d, Capacity: %d, Value: %v\n", output, len(output), cap(output), output)
		fmt.Printf("Expected >>> Type: %T, Length: %d, Capacity: %d, Value: %v\n", testCase.expected, len(testCase.expected), cap(testCase.expected), testCase.expected)

		if testCase.expected == nil || output == nil {
			t.Errorf("Expected output or actual output should not be nil")
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
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
