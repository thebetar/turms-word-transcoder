package transcoder

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	got := Tokenize("Hello, World!  It's  me.")
	want := []string{"hello", ",", "world", "!", "it's", "me", "."}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokenize() = %v, want %v", got, want)
	}
}

func TestBuild(t *testing.T) {
	tr := Build("the cat saw the dog")
	// <unk>, the, cat, saw, dog
	if tr.Size() != 5 {
		t.Fatalf("Size() = %d, want 5", tr.Size())
	}
}

func TestEncode(t *testing.T) {
	tr := Build("the cat saw the dog")
	got := tr.Encode("The dog saw a cat")
	want := []int{1, 4, 3, 0, 2} // "a" is unknown
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Encode() = %v, want %v", got, want)
	}
}

func TestDecode(t *testing.T) {
	tr := Build("hello, world!")
	got := tr.Decode([]int{1, 2, 3, 4, 99})
	want := "hello, world! <unk>"
	if got != want {
		t.Fatalf("Decode() = %q, want %q", got, want)
	}
}

func TestRoundTrip(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog."
	tr := Build(text)
	if got := tr.Decode(tr.Encode(text)); got != text {
		t.Fatalf("round trip = %q, want %q", got, text)
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "vocab.json")
	tr := Build("one two three")
	if err := tr.Save(path); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(loaded.Encode("three two one"), tr.Encode("three two one")) {
		t.Fatal("loaded vocabulary encodes differently than the original")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatal("expected an error for a missing file")
	}
}
