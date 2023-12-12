// highly inspired by jun06t
// cf. https://christina04.hatenablog.com/entry/routing-with-radix-tree
// cf. https://github.com/jun06t/go-sample/tree/master/radix-tree
// MIT License

package api

import (
	"fmt"
	"strings"
)

type Node struct {
	Part     string
	Children []*Node
	IsWild   bool
	Route    Route
}

func (n *Node) insert(parttern string, route Route) {
	parts := strings.Split(parttern, "/")[1:]

	for _, part := range parts {
		child := n.matchChild(part)
		if child == nil {
			fmt.Printf("part: %s\n", part)
			child = &Node{
				Part:   part,
				IsWild: isWild(part),
			}
			n.Children = append(n.Children, child)
		}
		n = child
	}
	n.Route = route
}

func (n *Node) matchChild(part string) *Node {
	for _, child := range n.Children {
		if child.Part == part || child.IsWild {
			return child
		}
	}
	return nil
}

func (n *Node) search(path string) Route {
	parts := strings.Split(path, "/")[1:]

	for _, part := range parts {
		child := n.matchChild(part)
		if child == nil {
			return Route{}
		}
		n = child
	}
	return n.Route
}

func isWild(part string) bool {
	if len(part) == 0 {
		return false
	}
	return part[0] == ':' || part[0] == '*' || (part[0] == '{' && part[len(part)-1] == '}')
}
