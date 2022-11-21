package wordy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// we're just going to make a hash map and hold all the
// words in there so that we don't need a real database

type dictResponse struct {
	Page  int                `json:"page"`
	RPP   int                `json:"rpp"`
	Total int                `json:"total"`
	Data  []dictWordResponse `json:"data"`
}

type dictWordResponse struct {
	ID   int    `json:"id"`
	Word string `json:"word"`
}

func GetWords() []string {
	path := os.Getenv("BASE_PATH") + "dictionary"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		panic(errors.New("error fetching dictionaries: " + err.Error()))
	}
	req.Header.Set("Authorization", "Basic dGVzdHVzZXI6dGVzdHBhc3M=")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(errors.New("error fetching dictionaries: " + err.Error()))
	}
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("error retreiving dictionary: expected 200, got %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	dr := dictResponse{}
	err = json.NewDecoder(resp.Body).Decode(&dr)
	if err != nil && err != io.EOF {
		panic(err)
	}

	pages := dr.Total / dr.RPP

	if dr.Total%dr.RPP != 0 {
		pages++
	}

	dict := []string{}
	for _, wr := range dr.Data {
		dict = append(dict, wr.Word)
	}

	// rinse and repeat
	for i := 2; i < pages; i++ {
		current := fmt.Sprintf("%s?page=%d", path, i)
		req, err := http.NewRequest("GET", current, nil)
		if err != nil {
			panic(errors.New("error fetching dictionaries: " + err.Error()))
		}
		req.Header.Set("Authorization", "Basic dGVzdHVzZXI6dGVzdHBhc3M=")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Connection", "keep-alive")
		resp, err := client.Do(req)
		if err != nil {
			panic(errors.New("error fetching dictionaries: " + err.Error()))
		}
		if resp.StatusCode != 200 {
			panic(fmt.Errorf("error retreiving dictionary: expected 200, got %d", resp.StatusCode))
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&dr)
		if err != nil && err != io.EOF {
			fmt.Println("error unmarshalling")
			panic(err)
		}

		for _, wr := range dr.Data {
			dict = append(dict, wr.Word)
		}
	}
	return dict
}
