package transcoder

import (
	"path/filepath"
	"reflect"
	"testing"
)

// The example from the video: troll2 and gymkata are used the same way.
const statquest = "Troll2 is great! Gymkata is great!"

func TestTokenize(t *testing.T) {
	got := Tokenize("Hello, World!  It's  me.")
	want := []string{"hello", ",", "world", "!", "it's", "me", "."}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokenize() = %v, want %v", got, want)
	}
}

func TestVocabulary(t *testing.T) {
	v := NewVocabulary(Tokenize("the cat saw the dog"))
	if v.Size() != 4 {
		t.Fatalf("Size() = %d, want 4", v.Size())
	}
	if id, ok := v.ID("Dog"); !ok || id != 3 || v.Word(id) != "dog" {
		t.Fatalf("ID(\"Dog\") = %d, %v", id, ok)
	}
	if _, ok := v.ID("bird"); ok {
		t.Fatal("bird should not be in the vocabulary")
	}
}

func TestTrainingLowersLoss(t *testing.T) {
	_, before := Train(statquest, 2, 1, 0.1, 1)
	_, after := Train(statquest, 2, 200, 0.1, 1)
	if after >= before {
		t.Fatalf("loss went from %.3f to %.3f, want it to go down", before, after)
	}
}

func TestSimilarWordsGetSimilarNumbers(t *testing.T) {
	tr, _ := Train(statquest, 2, 200, 0.1, 1)
	matches, ok := tr.Similar("troll2", 1)
	if !ok || matches[0].Word != "gymkata" {
		t.Fatalf("most similar to troll2 = %v, want gymkata", matches)
	}
}

func TestEncodeDecodePredictsNextWord(t *testing.T) {
	tr, _ := Train(statquest, 2, 200, 0.1, 1)
	for word, next := range map[string]string{"troll2": "is", "gymkata": "is", "is": "great", "great": "!"} {
		numbers, ok := tr.Encode(word)
		if !ok {
			t.Fatalf("Encode(%q) not found", word)
		}
		if got := tr.Decode(numbers); got != next {
			t.Errorf("Decode(Encode(%q)) = %q, want %q", word, got, next)
		}
	}
}

func TestUnknownWord(t *testing.T) {
	tr, _ := Train(statquest, 2, 1, 0.1, 1)
	if _, ok := tr.Encode("unknown"); ok {
		t.Fatal("Encode should fail for an unknown word")
	}
	if _, ok := tr.Similar("unknown", 3); ok {
		t.Fatal("Similar should fail for an unknown word")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "model.json")
	tr, _ := Train(statquest, 2, 50, 0.1, 1)
	if err := tr.Save(path); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	a, _ := tr.Encode("great")
	b, _ := loaded.Encode("great")
	if !reflect.DeepEqual(a, b) || loaded.Size() != tr.Size() {
		t.Fatal("loaded network differs from the saved one")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatal("expected an error for a missing file")
	}
}
