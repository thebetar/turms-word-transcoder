# Introduction

The main goal of this project is to explore how a word encoder and decoder work and how to build one in Go, as simply as possible.

The name `turms-word-transcoder` comes from Turms, the messenger of the gods in Etruscan mythology (https://en.wikipedia.org/wiki/Turms). A messenger takes a message, carries it to the other side and delivers it without changing what it means, which is exactly what an encoder and decoder do with text.

## How it works

Computers (and language models) work with numbers, not words. A word transcoder translates between the two in three steps:

1. **Tokenize**: the text is lowercased and split into words and punctuation marks. `"Hello, world!"` becomes `["hello", ",", "world", "!"]`.
2. **Build a vocabulary**: every unique token in a corpus (a text file to learn from) gets its own ID, in the order it first appears. ID `0` is reserved for `<unk>`, the unknown token.
3. **Encode and decode**: encoding looks up the ID of every token, decoding looks up the token for every ID. Words that are not in the vocabulary are encoded as `0` and come back as `<unk>`.

All of this lives in [`internal/transcoder/transcoder.go`](internal/transcoder/transcoder.go). The vocabulary is saved as a JSON list of words, where the position of a word is its ID.

## Usage

```sh
make build

# Build a vocabulary from the example corpus
./bin/turms-word-transcoder build data/corpus.txt vocab.json
# saved 43 words to vocab.json

# Text to IDs ("hermes" is not in the corpus, so it becomes 0)
./bin/turms-word-transcoder encode vocab.json "Hermes carried the message to the gods."
# 0 12 3 24 22 3 6 10

# IDs back to text
./bin/turms-word-transcoder decode vocab.json 0 12 3 24 22 3 6 10
# <unk> carried the message to the gods.
```

Useful make targets:

- `make run ARGS="..."` — run the CLI without building
- `make build` — build a binary into `bin/turms-word-transcoder`
- `make test` — run unit tests
- `make coverage` — run tests with a coverage report

Pushes to `master` run the unit tests with coverage in GitHub Actions.
