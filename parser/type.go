package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

var lexerRules = lexer.Rules{
	"Root": {
		{"Ident", `\w+`, nil},
		{"Whitespace", `\s+`, nil},
		{"Punct", `[\[\],;]`, nil},
		{"Equals", `=`, lexer.Push("Args")},
	},
	"Args": {
		{"Ident", `\w+`, nil},
		{"Quote", `'`, nil},
		{"Colon", `:`, nil},
		{"Equals", `=`, lexer.Push("SingleArg")},
		lexer.Return(),
	},
	"SingleArg": {
		{"ArgValue", `[^[\[\]=;,:']+`, lexer.Pop()},
	},
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

type SingleArg struct {
	Name  *string `(@Ident "=")?`
	Value *string `@ArgValue`
}

type Filter struct {
	InLinks  []string    `("[" @Ident "]")*`
	Filter   string      `@Ident`
	Args     []SingleArg `("=" ((@@ (":" @@)*) | (Quote? @@ (":" @@)* Quote?)))?`
	OutLinks []string    `("[" @Ident "]")*`
}

func (f Filter) String() string {
	s := f.Filter
	for _, input := range f.InLinks {
		s += "\n\tInput: " + input
	}
	for _, arg := range f.Args {
		s += "\n\tArg: " + *arg.Name + "=" + *arg.Value
	}
	for _, output := range f.OutLinks {
		s += "\n\tOutput: " + output
	}
	s += "\n"
	return s
}

type Filterchain struct {
	Filters []*Filter `@@ ("," Whitespace? @@)*`
}

func (fc Filterchain) String() string {
	s := ""
	for _, filter := range fc.Filters {
		s += filter.String()
	}
	return s
}

type Filtergraph struct {
	Filterchains []*Filterchain `@@ (";" Whitespace? @@)*`
}

func (fg Filtergraph) String() string {
	s := ""
	for _, chain := range fg.Filterchains {
		s += fmt.Sprintln(chain)
	}
	return s
}
