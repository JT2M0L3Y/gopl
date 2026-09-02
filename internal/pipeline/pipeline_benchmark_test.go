package pipeline

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func BenchmarkPipeline(b *testing.B) {
	for b.Loop() {
		var output bytes.Buffer
		if _, err := Run(strings.NewReader(smokeProgram), strings.NewReader(""), &output); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPipelineStages(b *testing.B) {
	var lex, parse, semantic, generate, execute time.Duration
	for b.Loop() {
		var output bytes.Buffer
		metrics, err := Run(strings.NewReader(smokeProgram), strings.NewReader(""), &output)
		if err != nil {
			b.Fatal(err)
		}
		lex += metrics.Lex
		parse += metrics.Parse
		semantic += metrics.Semantic
		generate += metrics.Generate
		execute += metrics.Execute
	}
	b.ReportMetric(float64(lex.Nanoseconds())/float64(b.N), "lex-ns/op")
	b.ReportMetric(float64(parse.Nanoseconds())/float64(b.N), "parse-ns/op")
	b.ReportMetric(float64(semantic.Nanoseconds())/float64(b.N), "semantic-ns/op")
	b.ReportMetric(float64(generate.Nanoseconds())/float64(b.N), "generate-ns/op")
	b.ReportMetric(float64(execute.Nanoseconds())/float64(b.N), "execute-ns/op")
}
