package wordy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var stateCorrect string = "correct"
var stateOut string = "not-in-word"
var stateIn string = "in-word"

var start = []string{"raise", "stare", "crane"}

type wordyPlayResponse struct {
	Data wordyPlayData `json:"data"`
}

type wordyPlayData struct {
	Id        int                   `json:"id"`
	State     string                `json:"state"`
	Word      string                `json:"word"`
	CorrectIn int                   `json:"correct_in"`
	Guesses   map[int]guessResponse `json:"guesses"`
	Error     string                `json:"error"`
}

type guessRequest struct {
	Guesses []string `json:"guesses"`
}

type guessResponse struct {
	Word  string       `json:"word"`
	Tiles map[int]tile `json:"tiles"`
}

type tile struct {
	Letter string `json:"letter"`
	State  string `json:"state"`
}

// PlayWordy takes an int argument representing the wordy
// puzzle number. It then plays that puzzle until it either
// solves it or loses
func PlayWordy(dict []string, v int) (string, bool) {
	// pick a starting word from our start list
	rand.Seed(time.Now().UnixNano())
	guesses := []string{start[rand.Intn(len(start)-1)]}
	// submit that word as our first try
	playResp, err := submit(guesses, v)
	if err != nil {
		panic(err)
	}
	// handle the response
	word := [5]string{}
	in := make(map[string][]int)
	out := make(map[string]bool)
	victory := evaluateResponse(guesses, &word, in, out, playResp)
	if victory == 1 {
		return "Victory! First Try!", true
	}
	// make a for loop that plays until success or defeat
	errs := 0
	tries := 1
	for victory == 0 && errs < 10 {
		// keep doing this
		var s string
		s, dict = chooseNext(word, dict, in, out)
		for _, word := range guesses {
			if s == word {
				// don't try the same word try
				continue
			}
		}
		guesses = append(guesses, s)
		fmt.Printf("submitting guess %d\n", len(guesses))
		playResp, err = submit(guesses, v)
		if err != nil {
			// just try again
			errs++
			continue
		}
		tries++

		victory = evaluateResponse(guesses, &word, in, out, playResp)
	}
	if victory == 1 {
		return fmt.Sprintf("Victory in %d guesses!", tries), true
	}
	return "defeat :(", false
}

func submit(guesses []string, v int) (*wordyPlayResponse, error) {
	// submit word to the api
	url := fmt.Sprintf("%swordy/%d/play", os.Getenv("BASE_PATH"), v)
	submission := guessRequest{Guesses: guesses}
	fmt.Println(submission)
	data, err := json.Marshal(submission)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", os.Getenv("AUTH"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("error submitting guess: expected 200, got %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	// unmarshall into struct
	wordResp := wordyPlayResponse{}
	err = json.NewDecoder(resp.Body).Decode(&wordResp)
	if err != nil {
		return nil, err
	}
	return &wordResp, nil
}

func evaluateResponse(guesses []string, word *[5]string, in map[string][]int, out map[string]bool, w *wordyPlayResponse) int {
	// uncomment to print response for debugging purposes
	// fmt.Println(w)
	// check to see if we won
	if w.Data.CorrectIn > 0 {
		return 1
	} else if len(guesses) > 5 {
		return -1
	}
	// evaluate letters that are correct
	// and add those letters to the word
	for i, tile := range w.Data.Guesses[len(w.Data.Guesses)-1].Tiles {
		if tile.State == stateCorrect {
			word[i] = tile.Letter
			// if the letter exists in "in", delete it
			delete(in, tile.Letter)
		}
	}
	// evaluate letters that are not in the word
	// add those letters to "out"
	for _, tile := range w.Data.Guesses[len(w.Data.Guesses)-1].Tiles {
		if tile.State == stateOut {
			out[tile.Letter] = true
		}
	}
	// evaluate letters that are in the word, but not in the correct position
	for i, tile := range w.Data.Guesses[len(w.Data.Guesses)-1].Tiles {
		if tile.State == stateIn {
			// if that letter exists in "in", append its most recent position
			// vscode will throw "unneccessary guard", but it's important that
			// we know whether to append or create the slice
			if _, ok := in[tile.Letter]; ok {
				in[tile.Letter] = append(in[tile.Letter], i)
			} else {
				// otherwise, add it to "in" and create a slice that contains its position
				in[tile.Letter] = []int{i}
			}
		}
	}
	return 0
}
