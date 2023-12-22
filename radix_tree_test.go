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
		pattern  string
		expected *Node
	}{
		{
			name:    "simple",
			pattern: "/",
			expected: &Node{
				Part: "",
				Children: []*Node{
					{
						Part: "GET",
						Children: []*Node{
							{
								Part:     "",
								Children: nil,
								IsWild:   false,
								Route:    Route{Method: http.MethodGet, Pattern: "/", HandlerFunc: dummyHandlerFunc},
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
			name:    "/foo",
			pattern: "/foo",
			expected: &Node{
				Part: "",
				Children: []*Node{
					{
						Part: "GET",
						Children: []*Node{
							{
								Part:     "foo",
								Children: nil,
								IsWild:   false,
								Route:    Route{Method: http.MethodGet, Pattern: "/foo", HandlerFunc: dummyHandlerFunc},
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
			name:    "/foo/bar",
			pattern: "/foo/bar",
			expected: &Node{
				Part: "",
				Children: []*Node{
					{
						Part: "GET",
						Children: []*Node{
							{
								Part: "foo",
								Children: []*Node{
									{
										Part:     "bar",
										Children: nil,
										IsWild:   false,
										Route:    Route{Method: http.MethodGet, Pattern: "/foo/bar", HandlerFunc: dummyHandlerFunc},
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
				IsWild: false,
				Route:  Route{},
			},
		},
		{
			name:    "/foo/:id",
			pattern: "/foo/:id",
			expected: &Node{
				Part: "",
				Children: []*Node{
					{
						Part: "GET",
						Children: []*Node{
							{
								Part: "foo",
								Children: []*Node{
									{
										Part:     ":id",
										Children: nil,
										IsWild:   true,
										Key:      "id",
										Route:    Route{Method: http.MethodGet, Pattern: "/foo/:id", HandlerFunc: dummyHandlerFunc},
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
				IsWild: false,
				Route:  Route{},
			},
		},
		{
			name:    "/foo/:name/:id",
			pattern: "/foo/:name/:id",
			expected: &Node{
				Part: "",
				Children: []*Node{
					{
						Part: "GET",
						Children: []*Node{
							{
								Part: "foo",
								Children: []*Node{
									{
										Part: ":name",
										Children: []*Node{
											{
												Part:     ":id",
												Children: nil,
												IsWild:   true,
												Key:      "id",
												Route:    Route{Method: http.MethodGet, Pattern: "/foo/:name/:id", HandlerFunc: dummyHandlerFunc},
											},
										},
										Key:    "name",
										IsWild: true,
										Route:  Route{},
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
				IsWild: false,
				Route:  Route{},
			},
		},
		// TODO add more cases
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			n := &Node{}
			n.insert(http.MethodGet, c.pattern, Route{
				Method:      http.MethodGet,
				Pattern:     c.pattern,
				HandlerFunc: dummyHandlerFunc,
			})
			if !deepEqualNode(t, n, c.expected) {
				t.Errorf("\ngot %#v, \nwant %#v", n, c.expected)
				t.Log("got: ==========================")
				printChildren(t, n)
				t.Log("want: ==========================")
				printChildren(t, c.expected)
			}
		})
	}
}

func TestNodeSearch(t *testing.T) {
	dummyHandlerFunc := func(w http.ResponseWriter, r *http.Request) {}

	type input struct {
		method  string
		path    string
		pattern string
	}
	cases := []struct {
		name     string
		input    input
		expected Route
	}{
		{
			name:     "GET /foo",
			input:    input{method: http.MethodGet, path: "/foo", pattern: "/foo"},
			expected: Route{Method: http.MethodGet, Pattern: "/foo", HandlerFunc: dummyHandlerFunc, PathParamMap: map[string]string{}},
		},
		{
			name:     "POST /foo",
			input:    input{method: http.MethodPost, path: "/foo", pattern: "/foo"},
			expected: Route{Method: http.MethodPost, Pattern: "/foo", HandlerFunc: dummyHandlerFunc, PathParamMap: map[string]string{}},
		},
		{
			name:     "GET /foo/bar",
			input:    input{method: http.MethodGet, path: "/foo/bar", pattern: "/foo/bar"},
			expected: Route{Method: http.MethodGet, Pattern: "/foo/bar", HandlerFunc: dummyHandlerFunc, PathParamMap: map[string]string{}},
		},
		{
			name:  "GET /foo/baz",
			input: input{method: http.MethodGet, path: "/foo/baz", pattern: "/foo/:id"},
			expected: Route{Method: http.MethodGet, Pattern: "/foo/:id", HandlerFunc: dummyHandlerFunc, PathParamMap: map[string]string{
				"id": "baz",
			}},
		},
		{
			name:  "GET /foo/baz/123",
			input: input{method: http.MethodGet, path: "/foo/baz/123", pattern: "/foo/:id/:name"},
			expected: Route{Method: http.MethodGet, Pattern: "/foo/:id/:name", HandlerFunc: dummyHandlerFunc, PathParamMap: map[string]string{
				"id":   "baz",
				"name": "123",
			}},
		},
		// NOTE: this case is not supported
		// LIMITATION: some pattern which has same prefix is not supported another path param name
		// {
		// 	name:  "GET /foo/baz/123",
		// 	input: input{method: http.MethodGet, path: "/foo/baz/123", pattern: "/foo/:name/:id"},
		// 	expected: Route{Method: http.MethodGet, Pattern: "/foo/:name/:id", HandlerFunc: dummyHandlerFunc, PathParamMap: map[string]string{
		// 		"name": "baz",
		// 		"id":   "123",
		// 	}},
		// },
	}

	n := &Node{}
	sm := &ServeMux{tree: n}
	// register all routes
	for _, c := range cases {
		sm.handle(c.input.method, c.input.pattern, dummyHandlerFunc)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := n.search(c.input.method, c.input.path)
			testEqual(t, c.expected.Pattern, r.Pattern)
			testEqual(t, c.expected.Method, r.Method)
			for k := range c.expected.PathParamMap {
				t.Logf("k: %s, v: %s", k, c.expected.PathParamMap[k])
				t.Logf("k: %s, v: %s", k, r.PathParamMap[k])
				testEqual(t, c.expected.PathParamMap[k], r.PathParamMap[k])
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

func deepEqualNode(t *testing.T, a, b *Node) bool {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		return true
	} else {
		if a.IsWild != b.IsWild ||
			a.Part != b.Part ||
			a.Key != b.Key ||
			a.Route.Method != b.Route.Method ||
			a.Route.Pattern != b.Route.Pattern {
			return false
		}
		if len(a.Route.PathParamMap) != len(b.Route.PathParamMap) {
			return false
		}
		for k := range a.Route.PathParamMap {
			if a.Route.PathParamMap[k] != b.Route.PathParamMap[k] {
				return false
			}
		}
		if len(a.Children) != len(b.Children) {
			return false
		}
		for i := range a.Children {
			// t.Logf("i: %d, a.Children[i]: %#v", i, a.Children[i])
			// t.Logf("i: %d, b.Children[i]: %#v", i, b.Children[i])
			if !deepEqualNode(t, a.Children[i], b.Children[i]) {
				return false
			}
		}
		return true
	}
}

func testEqual[T comparable](t *testing.T, a, b T) bool {
	t.Helper()
	result := reflect.DeepEqual(a, b)
	if !result {
		t.Errorf("\ngot %#v, \nwant %#v", a, b)
	}
	return result
}
