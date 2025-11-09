package content

import "testing"

func TestNormalizeWaybackSnapshotURLKeepsTimestamp(t *testing.T) {
	input := "https://web.archive.org/web/20231110091531/https://example.com/article"
	want := "https://web.archive.org/web/20231110091531id_/https://example.com/article"

	if got := normalizeWaybackSnapshotURL(input); got != want {
		t.Fatalf("normalizeWaybackSnapshotURL(%q) = %q, want %q", input, got, want)
	}
}

func TestNormalizeWaybackSnapshotURLPreservesRawInput(t *testing.T) {
	input := "https://web.archive.org/web/20231110091531id_/https://example.com/article"

	if got := normalizeWaybackSnapshotURL(input); got != input {
		t.Fatalf("normalizeWaybackSnapshotURL(%q) = %q, want the original value", input, got)
	}
}
