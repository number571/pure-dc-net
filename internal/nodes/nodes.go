package nodes

import "github.com/number571/pure-dc-net/pkg/dc"

func NodesKeysToGenerators(nodes Nodes) []dc.IGenerator {
	generators := make([]dc.IGenerator, 0, len(nodes))
	for _, n := range nodes {
		generators = append(generators, dc.NewHGenerator(n.SKey))
	}
	return generators
}
