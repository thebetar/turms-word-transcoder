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

The corpora in `data/` are religious texts from the Etruscans and their Roman neighbours, in English:

- [`liber-linteus.txt`](data/liber-linteus.txt): the [Liber Linteus](https://en.wikipedia.org/wiki/Liber_Linteus), the longest surviving Etruscan text, a ritual calendar of offerings to the gods. Etruscan is only partly understood, so this is the tentative translation of the passages that have one, from the notes on the [Wikipedia talk page](https://en.wikipedia.org/wiki/Talk:Liber_Linteus) ([CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/)). Words nobody can translate yet stay in Etruscan.
- [`tages.txt`](data/tages.txt): the myth of Tages, the child prophet who rose from a ploughed field and taught the Etruscans divination, from Cicero's _On Divination_ 2.50 ([Falconer, 1923](https://penelope.uchicago.edu/Thayer/E/Roman/Texts/Cicero/de_Divinatione/2*.html)) and Ovid's _Metamorphoses_ 15 ([Garth, Dryden et al., 1727](https://en.wikisource.org/wiki/Metamorphoses_(tr._Garth,_Dryden,_et_al.)/Book_XV)).
- [`roman-religion.txt`](data/roman-religion.txt): Roman prayers and rituals. Cato's farm offerings to Janus, Jupiter and Mars from _On Agriculture_ 132, 134 and 139-141 ([Hooper & Ash, 1934](https://penelope.uchicago.edu/Thayer/E/Roman/Texts/Cato/De_Agricultura/H*.html)), Cicero on Etruscan divination from _On Divination_ 1.92-93 ([Falconer, 1923](https://penelope.uchicago.edu/Thayer/E/Roman/Texts/Cicero/de_Divinatione/1*.html)), and Livy on how King Numa founded Roman religion from _History of Rome_ 1.19-21 ([Spillan, 1857](https://www.gutenberg.org/ebooks/19725)).

## Usage

```sh
make build

./bin/turms-word-transcoder train model.json data/*.txt
# trained on 1060 words, loss 3.826

# Word to numbers
./bin/turms-word-transcoder encode model.json decorated
# 1.474 -0.485 -1.287 1.831 0.125

# Numbers to the most likely next word
./bin/turms-word-transcoder decode model.json 1.474 -0.485 -1.287 1.831 0.125
# litter

# Words with similar numbers
./bin/turms-word-transcoder similar model.json wine
# pig          0.998
# barley       0.975
# around       0.970
# arms         0.957
# libation     0.957
```

With only about 3,800 words of text the results are noisy. Word2Vec is normally trained on millions of words.

Useful make targets:

- `make build` — build a binary into `bin/turms-word-transcoder`
- `make test` — run unit tests
- `make coverage` — run tests with a coverage report

Pushes to `master` run the unit tests with coverage in GitHub Actions.
