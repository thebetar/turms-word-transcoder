# Introduction

The main goal of this project is to explore how a word encoder and decoder work and how to build one in Go, as simply as possible.

The name `turms-word-transcoder` comes from Turms, the messenger of the gods in Etruscan mythology (https://en.wikipedia.org/wiki/Turms). A messenger takes a message, carries it to the other side and delivers it without changing what it means, which is exactly what an encoder and decoder do with text.

## How it works

Computers (and language models) work with numbers, not words. A word transcoder translates between the two in three steps:

1. **Tokenize**: the text is lowercased and split into words and punctuation marks. `"Hello, world!"` becomes `["hello", ",", "world", "!"]`.
2. **Build a vocabulary**: every unique token in a corpus (a text file to learn from) gets its own ID, in the order it first appears. ID `0` is reserved for `<unk>`, the unknown token.
3. **Encode and decode**: encoding looks up the ID of every token, decoding looks up the token for every ID. Words that are not in the vocabulary are encoded as `0` and come back as `<unk>`.

All of this lives in [`internal/transcoder/transcoder.go`](internal/transcoder/transcoder.go). The vocabulary is saved as a JSON list of words, where the position of a word is its ID.

## Data

Fitting the name, the example corpora in `data/` are Etruscan texts in English. No Etruscan myths survive in the Etruscan language itself: what is left is mostly short inscriptions, and their religious books (the _Etrusca Disciplina_) are lost. So there are two corpora:

- [`data/liber-linteus.txt`](data/liber-linteus.txt): the [Liber Linteus](https://en.wikipedia.org/wiki/Liber_Linteus), the longest surviving Etruscan text. It is a ritual calendar of offerings to the gods, written on linen that was later cut into strips to wrap an Egyptian mummy. Etruscan is still only partly understood, so no full translation exists. This file collects the tentative English translations of the passages that have one, taken from the column-by-column notes on the [Wikipedia talk page](https://en.wikipedia.org/wiki/Talk:Liber_Linteus) (based on L. B. van der Meer's 2007 edition, licensed [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/)). The Etruscan glosses and uncertain alternatives were removed, and words nobody can translate yet, such as `zarfneθ` or `cletram`, stay in Etruscan.
- [`data/tages.txt`](data/tages.txt): the myth of Tages, the child prophet who rose from a ploughed field and taught the Etruscans divination, as told by Roman writers in public domain English translations: Cicero, _On Divination_ 2.50 ([Falconer, 1923](https://penelope.uchicago.edu/Thayer/E/Roman/Texts/Cicero/de_Divinatione/2*.html)) and Ovid, _Metamorphoses_ 15 ([Garth, Dryden et al., 1727](https://en.wikisource.org/wiki/Metamorphoses_(tr._Garth,_Dryden,_et_al.)/Book_XV)).

## Usage

```sh
make build

# Build a vocabulary from the Liber Linteus
./bin/turms-word-transcoder build data/liber-linteus.txt linteus.json
# saved 302 words to linteus.json

# Text to IDs
./bin/turms-word-transcoder encode linteus.json "Pour a libation of new wine to Neptune."
# 168 64 38 79 69 57 20 150 12

# IDs back to text
./bin/turms-word-transcoder decode linteus.json 168 64 38 79 69 57 20 150 12
# pour a libation of new wine to neptune.

# Build a vocabulary from the Tages myth and try words that are not in it
./bin/turms-word-transcoder build data/tages.txt tages.json
# saved 168 words to tages.json

./bin/turms-word-transcoder encode tages.json "A boy sprang from the ploughed field and taught the Etruscans to read the gods."
# 8 43 28 105 1 18 15 24 163 1 106 31 0 1 0 33

./bin/turms-word-transcoder decode tages.json 8 43 28 105 1 18 15 24 163 1 106 31 0 1 0 33
# a boy sprang from the ploughed field and taught the etruscans to <unk> the <unk>.
```

The unknown words show the main weakness of a word-level transcoder: `gods` is unknown even though `god` might be in a bigger corpus, because every form of a word needs its own ID. Subword tokenizers such as byte pair encoding split words into smaller pieces to solve this.

Useful make targets:

- `make run ARGS="..."` — run the CLI without building
- `make build` — build a binary into `bin/turms-word-transcoder`
- `make test` — run unit tests
- `make coverage` — run tests with a coverage report

Pushes to `master` run the unit tests with coverage in GitHub Actions.
