package matcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	dburl := "http://" + os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT") + "/data"

	resp, err := http.Get(dburl)
	if err != nil {
		fmt.Print("Error getting data from server: ")
		fmt.Println(err)
		return &[]Matches{{Keywords: []string{""}, Phrase: "Derzeit besteht leider ein Problem mit der Datenbankverbindung. Bitte versuchen Sie es später erneut."}}
	}
	defer resp.Body.Close()
	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("Error reading response body: ")
		fmt.Println(err)
		return &[]Matches{{Keywords: []string{""}, Phrase: "Derzeit besteht leider ein Problem mit der Datenbankverbindung. Bitte versuchen Sie es später erneut."}}
	}

	var matches []Matches
	err = json.Unmarshal(byteValue, &matches)
	if err != nil {
		fmt.Print("Error unmarshalling JSON file: ")
		fmt.Println(err)
		return &[]Matches{{Keywords: []string{""}, Phrase: "Derzeit besteht leider ein Problem mit der Datenbankverbindung. Bitte versuchen Sie es später erneut."}}
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

type LogData struct {
	Tags    []string `json:"tags"`
	Problem string   `json:"problem"`
}

func logNoMatch(input string) {
	logData := LogData{Tags: []string{"unsolved"}, Problem: input}
	jsonPayload, err := json.Marshal(logData)

	if err != nil {
		fmt.Print("Error marshalling log data: ")
		fmt.Println(err)
	}

	url := "http://" + os.Getenv("UNSOLVEDHOST") + ":" + os.Getenv("UNSOLVEDPORT") + "/insert"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		fmt.Print("Error creating POST request: ")
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print("Error logging no match: ")
		fmt.Println(err)
	}
	defer resp.Body.Close()
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
		if len((*matches)[i].Keywords) > maxLength {
			maxLength = len((*matches)[i].Keywords)
			bestMatch = i
		}
		if len((*matches)[i].Keywords) == maxLength {
			if (*matches)[bestMatch].Keywords[0] == "" {
				bestMatch = i
			}
		}
	}

	if len(possibleMatches) == 1 && (*matches)[bestMatch].Keywords[0] == "" {
		fmt.Println("No match found")
		logNoMatch(input)
	}

	return (*matches)[bestMatch].Phrase
}
