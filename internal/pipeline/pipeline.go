package pipeline

import (
	"io"
	"time"

	"gopl/internal/ast"
	"gopl/internal/generator"
	"gopl/internal/lexer"
	"gopl/internal/parser"
	"gopl/internal/semantic"
	"gopl/internal/vm"
)

// Metrics records the time spent in each compiler stage.
type Metrics struct {
	Lex      time.Duration
	Parse    time.Duration
	Semantic time.Duration
	Generate time.Duration
	Execute  time.Duration
}

// Run executes the complete GoPL pipeline and returns stage timings.
func Run(source io.Reader, input io.Reader, output io.Writer) (Metrics, error) {
	var metrics Metrics

	start := time.Now()
	lex := lexer.New(source)

	start = time.Now()
	program, err := parser.New(lex).Parse()
	frontend := time.Since(start)
	metrics.Lex = lex.Duration()
	metrics.Parse = max(frontend-metrics.Lex, 0)
	if err != nil {
		return metrics, err
	}

	start = time.Now()
	if err := semantic.NewSemanticChecker().Check(program); err != nil {
		metrics.Semantic = time.Since(start)
		return metrics, err
	}
	metrics.Semantic = time.Since(start)

	start = time.Now()
	runtime := vm.New()
	runtime.SetInput(input)
	runtime.SetOutput(output)
	if err := generator.New(runtime).Generate(program); err != nil {
		metrics.Generate = time.Since(start)
		return metrics, err
	}
	metrics.Generate = time.Since(start)

	start = time.Now()
	err = runtime.Execute(false)
	metrics.Execute = time.Since(start)
	return metrics, err
}

// Parse is useful to tools and benchmarks that only need the front end.
func Parse(source io.Reader) (*ast.Program, error) {
	return parser.New(lexer.New(source)).Parse()
}
