package transcoder

import (
	"strings"
	"unicode"
)

// Tokenize splits lowercased text into words and punctuation marks.
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

// Vocabulary gives every unique word an ID, which is its input and output in the network.
type Vocabulary struct {
	words []string
	ids   map[string]int
}

func NewVocabulary(tokens []string) *Vocabulary {
	v := &Vocabulary{ids: map[string]int{}}
	for _, t := range tokens {
		if _, ok := v.ids[t]; !ok {
			v.ids[t] = len(v.words)
			v.words = append(v.words, t)
		}
	}
	return v
}

func (v *Vocabulary) ID(word string) (int, bool) {
	id, ok := v.ids[strings.ToLower(word)]
	return id, ok
}

func (v *Vocabulary) Word(id int) string { return v.words[id] }

func (v *Vocabulary) Size() int { return len(v.words) }
