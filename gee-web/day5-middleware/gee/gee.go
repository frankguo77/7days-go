package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type (
    RouterGroup struct{
		prefix			string
		middlewares		[]HandlerFunc
		parent			*RouterGroup
		engine 			*Engine
	}

	Engine struct{
		*RouterGroup
		router  *router
		groups  []*RouterGroup
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine, prefix:""}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefixp string) *RouterGroup {
	// log.Printf("GroupPrefix: %s", group.prefix)
	// log.Printf("Prefix: %s", prefixp)
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefixp,
		parent: group,
		engine: engine,
	}
	// log.Printf("newGroup.Prefix: %s", newGroup.prefix)
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
