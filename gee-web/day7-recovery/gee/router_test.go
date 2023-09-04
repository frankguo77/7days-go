package gee

import (
	// "fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestRouteParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/geekutu")

	if n == nil {
		t.Fatal("nil returned!")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("pattren mismatch!")
	}

	if ps["name"] != "geekutu" {
		t.Fatalf("Expected path[name] : %s, Got: %s", "geekutu", ps["name"])
	}

	// fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}