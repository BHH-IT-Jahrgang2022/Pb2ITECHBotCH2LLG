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
