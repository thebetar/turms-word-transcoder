package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"turms-word-transcoder/internal/transcoder"
)

const usage = `usage:
  turms-word-transcoder build  <corpus.txt> <vocab.json>   build a vocabulary from a text file
  turms-word-transcoder encode <vocab.json> <text...>      turn text into word IDs
  turms-word-transcoder decode <vocab.json> <ids...>       turn word IDs back into text`

func main() {
	if len(os.Args) < 4 {
		exit(usage)
	}

	command, file, args := os.Args[1], os.Args[2], os.Args[3:]

	switch command {
	case "build":
		build(file, args[0])
	case "encode":
		encode(file, strings.Join(args, " "))
	case "decode":
		decode(file, args)
	default:
		exit(usage)
	}
}

func build(corpusPath, vocabPath string) {
	corpus, err := os.ReadFile(corpusPath)
	if err != nil {
		exit(err.Error())
	}

	t := transcoder.Build(string(corpus))
	if err := t.Save(vocabPath); err != nil {
		exit(err.Error())
	}
	fmt.Printf("saved %d words to %s\n", t.Size(), vocabPath)
}

func encode(vocabPath, text string) {
	t := load(vocabPath)

	ids := t.Encode(text)
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.Itoa(id)
	}
	fmt.Println(strings.Join(parts, " "))
}

func decode(vocabPath string, args []string) {
	t := load(vocabPath)

	ids := make([]int, len(args))
	for i, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			exit(fmt.Sprintf("%q is not a valid ID", arg))
		}
		ids[i] = id
	}
	fmt.Println(t.Decode(ids))
}

func load(vocabPath string) *transcoder.Transcoder {
	t, err := transcoder.Load(vocabPath)
	if err != nil {
		exit(err.Error())
	}
	return t
}

func exit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
