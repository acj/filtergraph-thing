package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/pkg/errors"
)

func ParseFiltergraph(rawFiltergraph string) (*Filtergraph, error) {
	lexer, err := lexer.New(lexerRules)
	if err != nil {
		return nil, errors.Wrap(err, "creating lexer")
	}
	parser, err := participle.Build[Filtergraph](
		participle.Lexer(lexer),
		participle.UseLookahead(2),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "building parser")
	}

	return parser.ParseString("", rawFiltergraph)
}
