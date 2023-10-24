package main

import (
	"fmt"
	"io"

	graphviz "github.com/awalterschulze/gographviz"
	"github.com/pkg/errors"
)

//type Filter interface {
//	String() string
//}

//type Filterchain []Filter

type ScaleFilter struct {
	Width  int
	Height int
}

func (f *ScaleFilter) String() string {
	return fmt.Sprintf("scale=%d:%d", f.Width, f.Height)
}

func convertToFiltergraph(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", errors.Wrap(err, "couldn't read graph")
	}

	graphAst, err := graphviz.Parse(b)
	if err != nil {
		return "", errors.Wrap(err, "couldn't parse graph")
	}
	graph := &ChillGraph{graphviz.NewGraph()}
	if err := graphviz.Analyse(graphAst, graph); err != nil {
		return "", errors.Wrap(err, "couldn't analyze graph")
	}

	graph.Dump()

	return graph.ToFiltergraph()
}
