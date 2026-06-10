package ipfs

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUploadReturnsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	_, err := NewClient(server.URL).Upload([]byte("data"))
	if err == nil {
		t.Fatal("expected upload error")
	}
}

func TestLiveKuboUploadDownload(t *testing.T) {
	apiURL := os.Getenv("BIT_IPFS_API")
	if apiURL == "" {
		t.Skip("BIT_IPFS_API is not set")
	}

	client := NewClient(apiURL)
	payload := []byte("bit live ipfs verification\n")
	cid, err := client.Upload(payload)
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.Download(cid)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("downloaded payload mismatch: got %q, want %q", got, payload)
	}
}
