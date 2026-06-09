package cid

import "testing"

func TestCIDV0DigestRoundTrip(t *testing.T) {
	value := "QmVc67dGuQtDkHFbMDdn1WFUEDHn2RnAfF79byZYhhS51p"
	digest, err := DigestFromCIDV0(value)
	if err != nil {
		t.Fatal(err)
	}
	if got := CIDV0FromDigest(digest); got != value {
		t.Fatalf("got %s, want %s", got, value)
	}
}

func TestCIDV0DigestRejectsUnsupportedCID(t *testing.T) {
	if _, err := DigestFromCIDV0("bafybeigdyrzt5sfp7udm7hu76tfqwdcw7hry6qn2ip4wghc4ak2j6f4fvi"); err == nil {
		t.Fatal("expected CIDv1 to be rejected")
	}
}
