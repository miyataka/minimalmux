// highly inspired by jun06t
// cf. https://christina04.hatenablog.com/entry/routing-with-radix-tree
// cf. https://github.com/jun06t/go-sample/tree/master/radix-tree
// MIT License

package api

import (
	"strings"
)

type Node struct {
	Part     string
	Children []*Node
	IsWild   bool
	Key      string
	Route    Route
}

func (n *Node) insert(method, parttern string, route Route) {
	parts := strings.Split("/"+method+parttern, "/")[1:]

	for _, part := range parts {
		child := n.matchChild(part)
		if child == nil {
			child = &Node{
				Part:   part,
				IsWild: isWild(part),
				Key:    wildKey(part),
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

func (n *Node) search(method, path string) Route {
	parts := strings.Split("/"+method+path, "/")[1:]

	pMap := map[string]string{}
	for _, part := range parts {
		child := n.matchChild(part)
		if child == nil {
			return Route{}
		}
		// if part is wild, set path param
		if child.IsWild {
			pMap[child.Key] = part
		}
		n = child
	}
	n.Route.setPathParams(pMap)
	return n.Route
}

// isWild returns true if pattern part is a wild part
func isWild(part string) bool {
	if len(part) == 0 {
		return false
	}
	return part[0] == ':' || part[0] == '*' || (part[0] == '{' && part[len(part)-1] == '}')
}

// wildKey returns the keyword of the wild from pattern part
func wildKey(part string) string {
	if !isWild(part) {
		return ""
	}
	if part[0] == ':' || part[0] == '*' {
		return part[1:]
	}
	return part[1 : len(part)-1]
}
