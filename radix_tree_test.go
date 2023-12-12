package api

import (
	"reflect"
	"testing"
)

func TestNodeInsert(t *testing.T) {
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
						Part:     "",
						Children: nil,
						IsWild:   false,
						Route:    Route{},
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
						Part:     "foo",
						Children: nil,
						IsWild:   false,
						Route:    Route{},
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
			n.insert(c.path, Route{})
			if !reflect.DeepEqual(n, c.expected) {
				t.Errorf("\ngot %#v, \nwant %#v", n, c.expected)
			}
		})
	}
}
