# Introduction

The main goal of this project is to explore how a word encoder and decoder work and how to build one in Go, as simply as possible.

The name `turms-word-transcoder` comes from Turms, the messenger of the gods in Etruscan mythology (https://en.wikipedia.org/wiki/Turms). A messenger takes a message, carries it to the other side and delivers it without changing what it means, which is what an encoder and decoder do with words.

## How it works

This follows StatQuest's [Word Embedding and Word2Vec, Clearly Explained!!!](https://www.youtube.com/watch?v=viZrOnJclY0), which is where I learned how this works.

1. Every unique word in the corpus gets its own input and output in a small neural network.
2. Each input is connected to a few activation functions. The weights on those connections are the numbers that represent the word, its **embedding**. Going from a word to its numbers is the **encoder**.
3. The activation functions are connected to every output, and SoftMax turns the outputs into probabilities. Going from numbers back to the most likely word is the **decoder**.
4. The network is trained with backpropagation to predict the next word. Words that are used in the same way end up with similar numbers.

The code is in [`internal/transcoder`](internal/transcoder). Word2Vec adds Continuous Bag of Words, Skip-gram and Negative Sampling on top of this, which are left out to keep it simple.

## Data

The corpora in `data/` are Homer's epics in Samuel Butler's public domain prose translations from Project Gutenberg, about 270,000 words:

- [`iliad.txt`](data/iliad.txt): the _Iliad_, the siege of Troy ([Gutenberg #2199](https://www.gutenberg.org/ebooks/2199)).
- [`odyssey.txt`](data/odyssey.txt): the _Odyssey_, Ulysses' journey home ([Gutenberg #1727](https://www.gutenberg.org/ebooks/1727)).

Butler uses the Roman names for the gods, so Zeus is Jove and Odysseus is Ulysses.

## Usage

```sh
make build

# Takes a few minutes
./bin/turms-word-transcoder train model.json data/*.txt
# trained on 10133 words, loss 5.883

# Word to numbers
./bin/turms-word-transcoder encode model.json son
# -1.111 0.622 0.583 -1.120 -0.141 -0.140 -0.834 -0.338 -0.296 -0.815

# Numbers to the most likely next word
./bin/turms-word-transcoder decode model.json -1.111 0.622 0.583 -1.120 -0.141 -0.140 -0.834 -0.338 -0.296 -0.815
# of

# Words with similar numbers
./bin/turms-word-transcoder similar model.json achilles
# priam        0.977
# agamemnon    0.965
# patroclus    0.964
# menelaus     0.957
# hector       0.950
```

Useful make targets:

- `make build` — build a binary into `bin/turms-word-transcoder`
- `make test` — run unit tests
- `make coverage` — run tests with a coverage report

Pushes to `master` run the unit tests with coverage in GitHub Actions.
