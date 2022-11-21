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

	for i := 1; i < 20; i++ {
		fmt.Println()
		fmt.Println("****************************************")
		fmt.Printf("Playing wordy number %d\n", i)
		fmt.Println(wordy.PlayWordy(dict, i))
		fmt.Println("****************************************")
		fmt.Println()
	}
}
