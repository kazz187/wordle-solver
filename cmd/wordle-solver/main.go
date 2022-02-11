package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kazz187/wordle-solver/pkg/wordle"
)

const (
	wordsFilePath = "words.txt"
	wordLength    = 5
)

func main() {
	b, err := os.ReadFile(wordsFilePath)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to read file(%s): %w", wordsFilePath, err))
	}

	words := strings.Split(string(b), "\n")

	solver := wordle.NewSolver(words, wordLength)

	for {
		w := solver.Recommend()
		fmt.Println(w)
		var in string
		fmt.Scanln(&in)
		solver.FilterByResult(w, []wordle.Result(in))
	}
}
