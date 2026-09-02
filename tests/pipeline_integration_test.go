package tests

import (
	"bytes"
	"strings"
	"testing"

	"gopl/internal/pipeline"
)

func TestPipelineIntegration(t *testing.T) {
	program := `int add(int left, int right) {
  return left + right
}

void main() {
  print(add(20, 22))
}`

	var output bytes.Buffer
	_, err := pipeline.Run(strings.NewReader(program), strings.NewReader(""), &output)
	if err != nil {
		t.Fatal(err)
	}
	if output.String() != "42" {
		t.Fatalf("output = %q, want %q", output.String(), "42")
	}
}
