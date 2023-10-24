package main

import (
	"fmt"
	"io"
	"os"

	"github.com/acj/filtergraph-thing/parser"
)

func main() {
	path := os.Args[1]
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read graph at '%s': %s", path, err)
		os.Exit(1)
	}
	defer f.Close()

	//fg, err := convertToFiltergraph(f)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Couldn't convert graph to filtergraph: %s", err)
	//	os.Exit(1)
	//}
	//fmt.Println(fg)

	fgRaw, _ := io.ReadAll(f)
	f.Seek(0, 0)
	fg, err := parser.ParseFiltergraph(string(fgRaw))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse filtergraph: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", fg)
}
