package matcher

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Test() string {

	return "I'm alive!"
}

type Matches struct {
	Keywords []string `json:"keywords"`
	Phrase   string   `json:"answer"`
}

func LoadTable() *[]Matches {
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
