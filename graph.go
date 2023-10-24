package main

import (
	"fmt"
	"strings"

	graphviz "github.com/awalterschulze/gographviz"
)

type ChillGraph struct {
	*graphviz.Graph
}

func (g *ChillGraph) AddNode(parentGraph string, name string, attrs map[string]string) error {
	as := make(graphviz.Attrs)
	for k, v := range attrs {
		// Avoid the attribute validation in Attrs.Add
		as[graphviz.Attr(k)] = v
	}
	g.Nodes.Add(&graphviz.Node{Name: name, Attrs: as})
	g.Relations.Add(parentGraph, name)
	return nil
}

func (g *ChillGraph) ToFiltergraph() (string, error) {
	var chains []string
	for _, chain := range g.chains() {
		chains = append(chains, strings.Join(chain, ","))
	}

	return strings.Join(chains, ";"), nil
}

func (g *ChillGraph) rootNodes() (nodes []*graphviz.Node) {
	nodesWithInboundEdges := make(map[string]bool)

	for _, e := range g.Edges.Edges {
		nodesWithInboundEdges[e.Dst] = true
	}

	rootNodes := []*graphviz.Node{}
	for _, n := range g.Nodes.Nodes {
		if _, ok := nodesWithInboundEdges[n.Name]; !ok {
			rootNodes = append(rootNodes, n)
		}
	}
	return rootNodes
}

func (g *ChillGraph) adjacencyList() (adjList map[string][]string) {
	adjList = make(map[string][]string)
	for _, e := range g.Edges.Edges {
		adjList[e.Src] = append(adjList[e.Src], e.Dst)
	}
	return adjList
}

// TODO: return []Filterchain
func (g *ChillGraph) chains() [][]string {
	chains := [][]string{}
	for _, n := range g.rootNodes() {
		chains = append(chains, g.walk(n, make(map[string]bool)))
	}
	return chains
}

func (g *ChillGraph) walk(node *graphviz.Node, visited map[string]bool) []string {
	visited[node.Name] = true
	chain := []string{node.Name}
	for _, n := range g.adjacencyList()[node.Name] {
		if _, ok := visited[n]; ok {
			continue
		}
		chain = append(chain, g.walk(g.Nodes.Lookup[n], visited)...)
	}
	return chain
}

func (g *ChillGraph) Dump() {
	fmt.Printf("Edges: %d\n", len(g.Edges.Edges))
	fmt.Printf("Nodes: %d\n", len(g.Nodes.Nodes))

	for _, n := range g.Nodes.Nodes {
		fmt.Printf("Node: %s\n", n.Name)
	}

	for _, n := range g.rootNodes() {
		fmt.Printf("Root node: %s\n", n.Name)

		currentNode := n.Name
		for {
			adjList := g.adjacencyList()
			if len(adjList[currentNode]) == 0 {
				break
			}
			currentNode = adjList[currentNode][0]
			fmt.Printf(" -> %s\n", currentNode)
		}
	}

	for _, chain := range g.chains() {
		fmt.Printf("Chain: %s\n", chain)
	}
}
