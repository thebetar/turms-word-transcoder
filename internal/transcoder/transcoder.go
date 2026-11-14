// Package transcoder turns text into a list of word IDs and back again.
//
// It works in three steps:
//  1. Tokenize: split the text into words and punctuation marks.
//  2. Build: give every unique token in a corpus its own ID (its vocabulary).
//  3. Encode/Decode: look tokens up in the vocabulary to go from text to IDs and back.
package transcoder

import (
	"encoding/json"
	"os"
	"strings"
	"unicode"
)

// UnknownToken is used for words that are not in the vocabulary. It always has ID 0.
const UnknownToken = "<unk>"

// Transcoder holds a vocabulary in both directions.
type Transcoder struct {
	wordToID map[string]int
	idToWord []string
}

// Tokenize lowercases the text and splits it into words and punctuation marks.
// "Hello, world!" becomes ["hello", ",", "world", "!"].
func Tokenize(text string) []string {
	var tokens []string
	var word strings.Builder

	flush := func() {
		if word.Len() > 0 {
			tokens = append(tokens, word.String())
			word.Reset()
		}
	}

	for _, r := range strings.ToLower(text) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\'':
			word.WriteRune(r)
		case unicode.IsSpace(r):
			flush()
		default:
			flush()
			tokens = append(tokens, string(r))
		}
	}
	flush()

	return tokens
}

// FromWords creates a transcoder from a list of words, where the position of a
// word in the list is its ID. The unknown token is added at ID 0 when missing.
func FromWords(words []string) *Transcoder {
	t := &Transcoder{wordToID: map[string]int{}}
	t.add(UnknownToken)
	for _, w := range words {
		t.add(w)
	}
	return t
}

// Build creates a vocabulary from every unique token in the corpus, in the
// order they first appear.
func Build(corpus string) *Transcoder {
	return FromWords(Tokenize(corpus))
}

func (t *Transcoder) add(word string) {
	if _, ok := t.wordToID[word]; ok {
		return
	}
	t.wordToID[word] = len(t.idToWord)
	t.idToWord = append(t.idToWord, word)
}

// Size returns the number of tokens in the vocabulary, including the unknown token.
func (t *Transcoder) Size() int {
	return len(t.idToWord)
}

// Encode turns text into a list of IDs. Unknown words get ID 0.
func (t *Transcoder) Encode(text string) []int {
	tokens := Tokenize(text)
	ids := make([]int, len(tokens))
	for i, token := range tokens {
		ids[i] = t.wordToID[token] // a missing key returns 0, the unknown token
	}
	return ids
}

// Decode turns a list of IDs back into text. IDs outside the vocabulary
// become the unknown token. Punctuation is attached to the previous word.
func (t *Transcoder) Decode(ids []int) string {
	var out strings.Builder
	for i, id := range ids {
		word := UnknownToken
		if id >= 0 && id < len(t.idToWord) {
			word = t.idToWord[id]
		}
		if i > 0 && !isPunctuation(word) {
			out.WriteByte(' ')
		}
		out.WriteString(word)
	}
	return out.String()
}

func isPunctuation(word string) bool {
	r := []rune(word)
	return len(r) == 1 && unicode.IsPunct(r[0])
}

// Save writes the vocabulary to a JSON file as a list of words, where the
// position of a word is its ID.
func (t *Transcoder) Save(path string) error {
	data, err := json.MarshalIndent(t.idToWord, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a vocabulary file written by Save.
func Load(path string) (*Transcoder, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var words []string
	if err := json.Unmarshal(data, &words); err != nil {
		return nil, err
	}
	return FromWords(words), nil
}
