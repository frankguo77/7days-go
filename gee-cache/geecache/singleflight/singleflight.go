package singleflight

import (
	"sync"
)

type call struct {
	wg   sync.WaitGroup
	val  interface{}
	err  error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func()(interface{}, error)) (interface{}, error) {
	g.mu.Lock()
    if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// log.Println("i'm waiting")
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
    // log.Println("adding ", key)
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err

}