package main

import (
	"errors"
	"fmt"
	"wordy/internal/wordy"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(errors.New("unable to load environment variables: " + err.Error()))
	}
	fmt.Println("getting list of words")
	dict := wordy.GetWords()

	// for _, word := range dict {
	// 	if word == "midge" {
	// 		fmt.Println("Found midge")
	// 	}
	// }

	var victories, defeats int

	// there are 511 wordy puzzles, but 114 makes our app crash
	// due to the word not being in the dictionary list
	for i := 1; i < 511; i++ {
		// for i := 1; i < 114; i++ {
		fmt.Println()
		fmt.Println("****************************************")
		fmt.Printf("Playing wordy number %d\n", i)
		msg, victory := wordy.PlayWordy(dict, i)
		fmt.Println(msg)
		if victory {
			victories++
		} else {
			defeats++
		}
		fmt.Println("****************************************")
		fmt.Println()
	}

	fmt.Println()
	fmt.Printf("Total victories: %d\nTotal defeats: %d\n", victories, defeats)
}
