package wordy

import (
	"fmt"
	"math/rand"
	"strings"
)

func chooseNext(word [5]string, dict []string, in map[string][]int, out map[string]bool) (string, []string) {
	newDict := []string{}
	// iterate through the dictionary
	// for each word in the dictionary:
	for _, w := range dict {
		include := true
		// check that the letters in "word" match
		for i := 0; i < 5; i++ {
			// for each valid letter in word, if it's not matched in the
			// dictionary word, set include to false
			if word[i] != "" && string(w[i]) != word[i] {
				include = false
				continue
			}
		}
		// check that no letters from "out" are included
		for i := 0; i < 5; i++ {
			if _, ok := out[string(w[i])]; ok {
				include = false
				continue
			}
		}
		// check that each letter in "in" is included and check
		// that no letters from "in" are in invalid positions
		// for each letter in "in"
		for let, pos := range in {
			// check that the word contains the letter
			if !strings.Contains(w, let) {
				include = false
				continue
			}
			// check each index against each letter in the word
			for i := 0; i < 5; i++ {
				if string(w[i]) == let {
					for j := 0; j < len(pos); j++ {
						// if the current position is a position that is not allowed
						if i == pos[j] {
							include = false
							continue
						}
					}
				}
			}
		}

		// if we haven't ruled the word out, include it in the shortened word list
		if include {
			newDict = append(newDict, w)
		}
	}

	if len(newDict) < 1 {
		fmt.Println(dict)
		panic("Somehow don't have any words!")
	}

	pos := rand.Intn(len(newDict))
	return newDict[pos], newDict
}
