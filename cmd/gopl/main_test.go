package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"gopl/internal/pipeline"
)

func TestRunProgram(t *testing.T) {
	source, err := os.CreateTemp(t.TempDir(), "program-*.gopl")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := source.WriteString("void main() { print(2 + 3) }"); err != nil {
		t.Fatal(err)
	}
	if err := source.Close(); err != nil {
		t.Fatal(err)
	}

	contents, err := os.ReadFile(source.Name())
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if _, err := pipeline.Run(strings.NewReader(string(contents)), strings.NewReader(""), &output); err != nil {
		t.Fatal(err)
	}
	if output.String() != "5" {
		t.Fatalf("output = %q, want %q", output.String(), "5")
	}
}
