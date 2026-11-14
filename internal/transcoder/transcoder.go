// Package transcoder is the word embedding network from StatQuest's
// "Word Embedding and Word2Vec, Clearly Explained!!!":
// https://www.youtube.com/watch?v=viZrOnJclY0
package transcoder

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"
	"sort"
)

type Transcoder struct {
	vocab *Vocabulary
	in    [][]float64 // in[word] are the weights to the activation functions: the word's embedding
	out   [][]float64 // out[activation][word] are the weights to the outputs
}

// Train teaches the network to predict the next word in the corpus.
func Train(corpus string, dims, epochs int, rate float64, seed int64) (*Transcoder, float64) {
	tokens := Tokenize(corpus)
	vocab := NewVocabulary(tokens)
	rng := rand.New(rand.NewSource(seed))
	t := &Transcoder{vocab, randomMatrix(rng, vocab.Size(), dims), randomMatrix(rng, dims, vocab.Size())}

	var loss float64
	for e := 0; e < epochs; e++ {
		loss = 0
		for _, i := range rng.Perm(len(tokens) - 1) {
			word, _ := vocab.ID(tokens[i])
			next, _ := vocab.ID(tokens[i+1])
			loss += t.step(word, next, rate)
		}
		loss /= float64(len(tokens) - 1)
	}
	return t, loss
}

// step does one round of backpropagation and returns the cross entropy loss.
func (t *Transcoder) step(word, next int, rate float64) float64 {
	hidden := t.in[word]
	probs := t.softmax(hidden)
	loss := -math.Log(probs[next])

	probs[next] -= 1 // derivative of cross entropy with softmax
	hiddenGrad := make([]float64, len(hidden))
	for k := range t.out {
		for j, g := range probs {
			hiddenGrad[k] += t.out[k][j] * g
			t.out[k][j] -= rate * hidden[k] * g
		}
	}
	for k := range hidden {
		hidden[k] -= rate * hiddenGrad[k]
	}
	return loss
}

func (t *Transcoder) softmax(hidden []float64) []float64 {
	scores := make([]float64, t.vocab.Size())
	for k, h := range hidden {
		for j, w := range t.out[k] {
			scores[j] += h * w
		}
	}
	max, sum := scores[0], 0.0
	for _, s := range scores {
		max = math.Max(max, s)
	}
	for j := range scores {
		scores[j] = math.Exp(scores[j] - max)
		sum += scores[j]
	}
	for j := range scores {
		scores[j] /= sum
	}
	return scores
}

// Encode turns a word into its embedding.
func (t *Transcoder) Encode(word string) ([]float64, bool) {
	id, ok := t.vocab.ID(word)
	if !ok {
		return nil, false
	}
	return append([]float64(nil), t.in[id]...), true
}

// Decode turns an embedding into the most likely next word.
func (t *Transcoder) Decode(embedding []float64) string {
	probs := t.softmax(embedding)
	best := 0
	for j := range probs {
		if probs[j] > probs[best] {
			best = j
		}
	}
	return t.vocab.Word(best)
}

type Match struct {
	Word       string
	Similarity float64
}

// Similar returns the n words with the most similar embeddings.
func (t *Transcoder) Similar(word string, n int) ([]Match, bool) {
	id, ok := t.vocab.ID(word)
	if !ok {
		return nil, false
	}
	var matches []Match
	for j := range t.in {
		if j != id {
			matches = append(matches, Match{t.vocab.Word(j), cosine(t.in[id], t.in[j])})
		}
	}
	sort.Slice(matches, func(a, b int) bool { return matches[a].Similarity > matches[b].Similarity })
	return matches[:min(n, len(matches))], true
}

func (t *Transcoder) Size() int { return t.vocab.Size() }

func cosine(a, b []float64) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	return dot / math.Sqrt(na*nb)
}

func randomMatrix(rng *rand.Rand, rows, cols int) [][]float64 {
	m := make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		for j := range m[i] {
			m[i][j] = rng.Float64() - 0.5
		}
	}
	return m
}

type file struct {
	Words []string
	In    [][]float64
	Out   [][]float64
}

func (t *Transcoder) Save(path string) error {
	data, err := json.Marshal(file{t.vocab.words, t.in, t.out})
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func Load(path string) (*Transcoder, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f file
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return &Transcoder{NewVocabulary(f.Words), f.In, f.Out}, nil
}
