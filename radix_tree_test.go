package api

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNodeInsert(t *testing.T) {
	dummyHandlerFunc := func(w http.ResponseWriter, r *http.Request) {}
	cases := []struct {
		name     string
		path     string
		expected *Node
	}{
		{
			name: "simple",
			path: "/",
			expected: &Node{
				Part: "",
				Children: []*Node{
					&Node{
						Part: "GET",
						Children: []*Node{
							&Node{
								Part:     "",
								Children: nil,
								IsWild:   false,
								Route:    Route{Method: http.MethodGet, Path: "/", HandlerFunc: dummyHandlerFunc},
							},
						},
						IsWild: false,
						Route:  Route{},
					},
				},
				IsWild: false,
				Route:  Route{},
			},
		},
		{
			name: "/foo",
			path: "/foo",
			expected: &Node{
				Part: "",
				Children: []*Node{
					&Node{
						Part: "GET",
						Children: []*Node{
							&Node{
								Part:     "foo",
								Children: nil,
								IsWild:   false,
								Route:    Route{Method: http.MethodGet, Path: "/foo", HandlerFunc: dummyHandlerFunc},
							},
						},
						IsWild: false,
						Route:  Route{},
					},
				},
				IsWild: false,
				Route:  Route{},
			},
		},
		// TODO add more cases
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			n := &Node{}
			n.insert(http.MethodGet, c.path, Route{
				Method:      http.MethodGet,
				Path:        c.path,
				HandlerFunc: dummyHandlerFunc,
			})
			if !deepEqualExceptPointer(t, n, c.expected) {
				t.Errorf("\ngot %#v, \nwant %#v", n, c.expected)
				t.Log("got: ==========================")
				printChildren(t, n)
				t.Log("want: ==========================")
				printChildren(t, c.expected)
			}
		})
	}
}

func printChildren(t *testing.T, n *Node) {
	t.Helper()
	t.Logf("%#v", n)
	for _, v := range n.Children {
		printChildren(t, v)
	}
}

func deepEqualExceptPointer(t *testing.T, a, b *Node) bool {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		return true
	} else {
		if a.IsWild == b.IsWild && a.Part == b.Part {
			if a.Route.Method != b.Route.Method || a.Route.Path != b.Route.Path {
				return false
			}
		}
		for i := range a.Children {
			if !deepEqualExceptPointer(t, a.Children[i], b.Children[i]) {
				return false
			}
		}
		return true
	}
}
