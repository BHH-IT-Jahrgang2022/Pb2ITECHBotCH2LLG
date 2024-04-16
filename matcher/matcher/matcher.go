package matcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func Test() string {

	return "I'm alive!"
}

type Matches struct {
	Keywords []string `json:"keywords"`
	Phrase   string   `json:"response"`
}

func LoadTable() *[]Matches {
	// TODO: Update to DB connection instead of local JSON file

	/*
		jsonFile, err := os.Open("./data/data.json")
		if err != nil {
			fmt.Printf("Error opening JSON file: %v\n", err)
			return &[]Matches{}
		}
		defer jsonFile.Close()
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			fmt.Print("Error reading JSON file: ")
			fmt.Println(err)
			return &[]Matches{}
		}

		var matches []Matches
		err = json.Unmarshal(byteValue, &matches)
		if err != nil {
			fmt.Print("Error unmarshalling JSON file: ")
			fmt.Println(err)
			return &[]Matches{}
		}
	*/

	resp, err := http.Get("http://127.0.0.1:8080/data")
	if err != nil {
		fmt.Print("Error getting data from server: ")
		fmt.Println(err)
		return &[]Matches{}
	}
	defer resp.Body.Close()
	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("Error reading response body: ")
		fmt.Println(err)
		return &[]Matches{}
	}

	var matches []Matches
	err = json.Unmarshal(byteValue, &matches)
	if err != nil {
		fmt.Print("Error unmarshalling JSON file: ")
		fmt.Println(err)
		return &[]Matches{}
	}

	fmt.Println("Succesfully initialized the matcher with following matches:")
	printAllMatches(&matches)
	fmt.Println()

	return &matches
}

func printAllMatches(matches *[]Matches) {
	for _, match := range *matches {
		fmt.Println(match.Keywords)
		fmt.Println(match.Phrase)
	}
}

func Match(input string, matches *[]Matches) string {
	var possibleMatches []int
	for i, match := range *matches {
		var matchCounter int
		for _, keyword := range match.Keywords {
			r := regexp.MustCompile("(?i)" + keyword)
			if r.MatchString(input) {
				matchCounter++
				fmt.Println("Matched: ", keyword)
			}

		}
		if matchCounter == len(match.Keywords) {
			possibleMatches = append(possibleMatches, i)
		}
	}

	var maxLength int
	var bestMatch int

	for _, i := range possibleMatches {
		if len((*matches)[i].Keywords) >= maxLength {
			maxLength = len((*matches)[i].Keywords)
			bestMatch = i
		}
	}

	return (*matches)[bestMatch].Phrase
}
