package wordle

import (
	"fmt"
	"strings"
)

type Result rune

const (
	ResultRightSpot = Result('o')
	ResultWrongSpot = Result('x')
	ResultNoSpot    = Result('.')
)

type words map[string]bool

type tree []map[rune]words

func newTree(length int) tree {
	tree := make([]map[rune]words, length)
	for i := range tree {
		tree[i] = map[rune]words{}
	}
	return tree
}

func (t *tree) registerWord(word string) {
	for i, r := range []rune(word) {
		w, ok := (*t)[i][r]
		if ok {
			w[word] = true
		} else {
			(*t)[i][r] = words{word: true}
		}
	}
}

func (t *tree) unregisterWord(word string) {
	for i, r := range []rune(word) {
		w, ok := (*t)[i][r]
		if ok {
			delete(w, word)
			if len(w) == 0 {
				delete((*t)[i], r)
			}
		}
	}
}

func (t *tree) proveRightSpot(rune rune, index int) {
	// index に rune がない単語は正解から外す
	for r, w := range (*t)[index] {
		if r != rune {
			for word, _ := range w {
				t.unregisterWord(word)
			}
		}
	}
}

func (t *tree) proveWrongSpot(rune rune, index int) {
	for r, w := range (*t)[index] {
		// index に rune がある単語は正解から外す
		if r == rune {
			for word, _ := range w {
				t.unregisterWord(word)
			}
		} else {
			// rune を含まない単語は正解から外す
			for word, _ := range w {
				if !strings.Contains(word, string(rune)) {
					t.unregisterWord(word)
				}
			}
		}
	}
}

func (t *tree) proveNoSpot(rune rune) {
	// rune を含む単語は正解から外す
	for _, w := range (*t)[0] {
		for word, _ := range w {
			if strings.ContainsRune(word, rune) {
				t.unregisterWord(word)
			}
		}
	}
}

func (t *tree) recommend() string {
	words := words{}
	for _, w := range (*t)[0] {
		for word, _ := range w {
			words[word] = true
		}
	}

	for w, _ := range words {
		return w
	}
	return ""
}

type Solver struct {
	length      int
	tree        tree
	triedRunes  map[rune]bool
	foundRunes  map[rune]bool
	filterWords words
}

func NewSolver(words []string, length int) *Solver {
	s := &Solver{
		length:      length,
		tree:        newTree(length),
		triedRunes:  map[rune]bool{},
		foundRunes:  map[rune]bool{},
		filterWords: map[string]bool{},
	}
	for _, word := range words {
		s.tree.registerWord(word)
		if uniqRunes(word) {
			s.filterWords[word] = true
		}
	}
	return s
}

func (s *Solver) FilterByResult(word string, result []Result) {
	for i, r := range []rune(word) {
		s.triedRunes[r] = true

		for w, _ := range s.filterWords {
			if strings.ContainsRune(w, r) {
				delete(s.filterWords, w)
			}
		}

		switch result[i] {
		case ResultRightSpot:
			s.foundRunes[r] = true
			s.tree.proveRightSpot(r, i)
		case ResultWrongSpot:
			s.foundRunes[r] = true
			s.tree.proveWrongSpot(r, i)
		case ResultNoSpot:
			s.tree.proveNoSpot(r)
		}
	}
}

func (s *Solver) Recommend() string {
	if len(s.foundRunes) < (s.length/2 + 1) {
		fmt.Println("phase 1")
		for w, _ := range s.filterWords {
			return w
		}
	}
	fmt.Println("phase 2")
	return s.tree.recommend()
}

func uniqRunes(word string) bool {
	for i := 0; i < len(word)-1; i++ {
		for j := i + 1; j < len(word); j++ {
			if []rune(word)[i] == []rune(word)[j] {
				return false
			}
		}
	}
	return true
}
