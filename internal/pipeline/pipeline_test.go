package pipeline

import (
	"bytes"
	"strings"
	"testing"
)

const smokeProgram = `void main() {
  print(2 + 3)
}
`

func TestRun(t *testing.T) {
	var output bytes.Buffer
	_, err := Run(strings.NewReader(smokeProgram), strings.NewReader(""), &output)
	if err != nil {
		t.Fatal(err)
	}
	if output.String() != "5" {
		t.Fatalf("output = %q, want %q", output.String(), "5")
	}
}
