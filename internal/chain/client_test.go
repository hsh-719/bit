package chain

import "testing"

func TestGitHashBytes20RoundTrip(t *testing.T) {
	hash := "0123456789abcdef0123456789abcdef01234567"
	asBytes, err := GitHashToBytes20(hash)
	if err != nil {
		t.Fatal(err)
	}
	if got := Bytes20ToGitHash(asBytes); got != hash {
		t.Fatalf("got %s, want %s", got, hash)
	}
	if IsZeroBytes20(asBytes) {
		t.Fatal("non-zero hash reported as zero")
	}
	if !IsZeroBytes20([20]byte{}) {
		t.Fatal("zero hash reported as non-zero")
	}
}

func TestGitHashBytes20RejectsNonSHA1(t *testing.T) {
	if _, err := GitHashToBytes20("abcd"); err == nil {
		t.Fatal("expected short hash to fail")
	}
}
