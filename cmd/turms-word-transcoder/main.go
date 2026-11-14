package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"turms-word-transcoder/internal/transcoder"
)

const (
	dims   = 10 // activation functions per word, the video uses 2
	epochs = 1
	rate   = 0.05
	seed   = 1
)

const usage = `usage:
  turms-word-transcoder train   <model.json> <corpus.txt...>
  turms-word-transcoder encode  <model.json> <word>
  turms-word-transcoder decode  <model.json> <numbers...>
  turms-word-transcoder similar <model.json> <word>`

func main() {
	if len(os.Args) < 4 {
		exit(usage)
	}
	command, model, args := os.Args[1], os.Args[2], os.Args[3:]

	switch command {
	case "train":
		train(model, args)
	case "encode":
		encode(load(model), args[0])
	case "decode":
		decode(load(model), args)
	case "similar":
		similar(load(model), args[0])
	default:
		exit(usage)
	}
}

func train(model string, files []string) {
	var corpus strings.Builder
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			exit(err.Error())
		}
		corpus.Write(data)
		corpus.WriteString("\n")
	}
	t, loss := transcoder.Train(corpus.String(), dims, epochs, rate, seed)
	if err := t.Save(model); err != nil {
		exit(err.Error())
	}
	fmt.Printf("trained on %d words, loss %.3f\n", t.Size(), loss)
}

func encode(t *transcoder.Transcoder, word string) {
	numbers, ok := t.Encode(word)
	if !ok {
		exit(word + " is not in the corpus")
	}
	for _, n := range numbers {
		fmt.Printf("%.3f ", n)
	}
	fmt.Println()
}

func decode(t *transcoder.Transcoder, args []string) {
	if len(args) != dims {
		exit(fmt.Sprintf("expected %d numbers", dims))
	}
	numbers := make([]float64, dims)
	for i, arg := range args {
		n, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			exit(arg + " is not a number")
		}
		numbers[i] = n
	}
	fmt.Println(t.Decode(numbers))
}

func similar(t *transcoder.Transcoder, word string) {
	matches, ok := t.Similar(word, 5)
	if !ok {
		exit(word + " is not in the corpus")
	}
	for _, m := range matches {
		fmt.Printf("%-12s %.3f\n", m.Word, m.Similarity)
	}
}

func load(model string) *transcoder.Transcoder {
	t, err := transcoder.Load(model)
	if err != nil {
		exit(err.Error())
	}
	return t
}

func exit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
