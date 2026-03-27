package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_MissingCommand_ShowsHelpfulError(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run(nil, &out, &errOut)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(errOut.String(), "missing command") {
		t.Fatalf("expected missing command error, got %q", errOut.String())
	}
	if !strings.Contains(errOut.String(), "Usage:") {
		t.Fatalf("expected usage text, got %q", errOut.String())
	}
}

func TestRun_Done_MissingID_ShowsHelpfulError(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"done"}, &out, &errOut)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(errOut.String(), "missing id") {
		t.Fatalf("expected missing id error, got %q", errOut.String())
	}
	if !strings.Contains(errOut.String(), "usage: tm done <id>") {
		t.Fatalf("expected usage hint, got %q", errOut.String())
	}
}

func TestRun_Done_InvalidID_ShowsHelpfulError(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"done", "nope"}, &out, &errOut)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(errOut.String(), "invalid id") {
		t.Fatalf("expected invalid id error, got %q", errOut.String())
	}
	if !strings.Contains(errOut.String(), "UUID") {
		t.Fatalf("expected UUID hint, got %q", errOut.String())
	}
}
