package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		log.Printf("Start Log!")
		c.Next()
		log.Printf("End Log! - [%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}