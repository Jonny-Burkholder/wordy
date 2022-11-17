package wordy

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var stateCorrect string = "correct"
var stateOut string = "incorrect"
var stateIn string = "in_word"

var letters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

type dictResponse struct {
	Total int `json:"total"`
}

type dictWordResponse struct {
	Data wordData `json:"data"`
}

type wordData struct {
	Word string `json:"word"`
}

type wordyPlayResponse struct {
	Data  wordyData `json:"data"`
	Error string    `json:"error"`
}

type wordyData struct {
	CorrectIn int             `json:"correct_in"`
	Guesses   []guessResponse `json:"guesses"`
}

type guessResponse struct {
	Tiles []tile `json:"tiles"`
}

type tile struct {
	Letter string `json:"letter"`
	State  string `json:"state"`
}

// PlayWordy takes an int argument representing the wordy
// puzzle number. It then plays that puzzle until it either
// solves it or loses
func PlayWordy(v int) string {
	// get request to dictionary
	resp, err := http.Get(os.Getenv("BASE_PATH") + "/dictionary")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dR := dictResponse{}
	err = json.NewDecoder(resp.Body).Decode(&dR)
	if err != nil {
		panic(err)
	}
	// pick a random number in dictionary range
	rand.Seed(time.Now().UnixMicro())
	num := rand.Intn(dR.Total)

	// get that word as our starting word
	resp, err = http.Get(fmt.Sprintf("%s/dictionary%d", os.Getenv("BASE_PATH"), num))
	if err != nil {
		panic(err)
	}

	wR := dictWordResponse{}
	// not actually sure if it's necessary to declare the decoder a second time
	err = json.NewDecoder(resp.Body).Decode(&wR)
	if err != nil {
		panic(err)
	}
	guesses := []string{wR.Data.Word}
	// submit that word as our first try
	playResp, err := submit(guesses, v)
	if err != nil {
		panic(err)
	}
	// handle the response
	word := [5]string{}
	in := make(map[string][]int)
	out := make(map[string]bool)
	victory := evaluateResponse(guesses, word, in, out, playResp)
	if victory == 1 {
		return "victory!"
	}
	// make a for loop that plays until success or defeat
	errs := 0
	for victory == 0 && errs < 10 {
		// keep doing this
		s := shuffle(word, in, out)
		for _, word := range guesses {
			if s == word {
				// don't try the same word try
				continue
			}
		}
		if !isWord(s) {
			continue
		}
		guesses = append(guesses, s)
		playResp, err = submit(guesses, v)
		if err != nil {
			// just try again
			continue
		}
		victory = evaluateResponse(guesses, word, in, out, playResp)
	}
	if victory == 1 {
		return "victory!"
	}
	return "defeat :("
}

// isWord checks to see if a word is, you know, a word
func isWord(s string) bool {
	// send request to dictionary with the string
	client := http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("BASE_PATH")+"/"+s, nil)
	if err != nil {
		return false
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return true
}

func shuffle(word [5]string, in map[string][]int, out map[string]bool) string {
	try := word
	// shuffle letters that we know are in the word
	for let, vals := range in {
		// generate a random number between 0 and 4 that is not already taken
		pos := rand.Intn(5)
		// check to see if i is an incorrect position
		for _, val := range vals {
			if val == pos {
				// do nothing
			} else {
				if try[pos] == "" {
					try[pos] = let
					break
				}
			}
		}
		var empty int
		// check to see if the word is complete
		empty = 0
		for _, let := range try {
			if let == "" {
				empty++
			}
		}
		for empty > 0 {
			for _, let := range letters {
				// check to make sure the letter is not out
				if _, ok := out[let]; ok != true {
					// randomly place the letter in the word
					for {
						pos := rand.Intn(5)
						if try[pos] == "" {
							try[pos] = let
						}
					}
				}
			}
			// check to see if the word is complete
			empty = 0
			for _, let := range try {
				if let == "" {
					empty++
				}
			}
		}
	}

	// convert try into a string
	s := ""
	for i := 0; i < 5; i++ {
		s += try[i]
	}
	return s
}

func submit(guesses []string, v int) (*wordyPlayResponse, error) {
	// submit word to the api
	url := fmt.Sprintf("%s/wordy/%d/play", os.Getenv("BASE_URL"), v)
	s := strings.Join(guesses, ",")
	req, err := http.NewRequest("GET", url, strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
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

func evaluateResponse(guesses []string, word [5]string, in map[string][]int, out map[string]bool, w *wordyPlayResponse) int {
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
			if _, ok := in[tile.Letter]; ok != false {
				delete(in, tile.Letter)
			}
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
			if _, ok := in[tile.Letter]; ok != false {
				in[tile.Letter] = append(in[tile.Letter], i)
			} else {
				// otherwise, add it to "in" and create a slice that contains its position
				in[tile.Letter] = []int{i}
			}
		}
	}
	return 0
}
