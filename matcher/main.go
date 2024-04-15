package main

import (
	"fmt"
	"matcher/matcher"
)

func main() {
	fmt.Println(matcher.Test())
	matches := matcher.LoadTable()
	fmt.Println((*matches)[0].Keywords)
}
