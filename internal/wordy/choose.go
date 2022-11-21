package wordy

import (
	"fmt"
	"math/rand"
	"strings"
)

func chooseNext(word [5]string, dict []string, in map[string][]int, out map[string]bool) (string, []string) {
	fmt.Println(len(dict))
	if len(dict) < 100 {
		fmt.Println(dict)
	}
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
				if len(dict) < 100 {
					fmt.Printf("%s does not match pattern %s", w, word)
				}
			}
		}
		// check that no letters from "out" are included
		for i := 0; i < 5; i++ {
			if _, ok := out[string(w[i])]; ok {
				include = false
				if len(dict) < 100 {
					fmt.Printf("%s contains the unallowed letter %s\n", w, string(w[i]))
				}
			}
		}
		// check that each letter in "in" is included and check
		// that no letters from "in" are in invalid positions
		// for each letter in "in"
		for let, pos := range in {
			// check that the word contains the letter
			if !strings.Contains(w, let) {
				include = false
				if len(dict) < 100 {
					fmt.Printf("%s does not contain the necessary letter %s\n", w, let)
				}
				continue
			}
			// check each index against each letter in the word
			for i := 0; i < 5; i++ {
				if string(w[i]) == let {
					for j := 0; j < len(pos); j++ {
						// if the current position is a position that is not allowed
						if i == j {
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
		} else {
			if len(dict) > 100 {
				fmt.Printf("excluding %s\n", w)
			}
		}
	}

	if len(newDict) < 1 {
		panic("Somehow don't have any words!")
	}

	pos := rand.Intn(len(newDict))
	return newDict[pos], newDict
}
