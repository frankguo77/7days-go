package main

/*
(1)
$ curl -i http://localhost:9999/
HTTP/1.1 200 OK
Date: Mon, 12 Aug 2019 16:52:52 GMT
Content-Length: 18
Content-Type: text/html; charset=utf-8
<h1>Hello Gee</h1>

(2)
$ curl "http://localhost:9999/hello?name=geektutu"
hello geektutu, you're at /hello

(3)
$ curl "http://localhost:9999/hello/geektutu"
hello geektutu, you're at /hello/geektutu

(4)
$ curl "http://localhost:9999/assets/css/geektutu.css"
{"filepath":"css/geektutu.css"}

(5)
$ curl "http://localhost:9999/xxx"
404 NOT FOUND: /xxx
*/

import (
	"gee"
	"net/http"
	"log"
	// "time"
)

func onlyForV21() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		// t := time.Now()
		log.Printf("group v21 - Start")
		// if a server error occurred
		c.Next()
		// Calculate resolution time
		log.Printf("group v21 - End")
	}
}

func onlyForV22() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		// t := time.Now()
		log.Printf("group v22 - Satrt")
		// if a server error occurred
		// c.Next()
		// Calculate resolution time
		log.Printf("group v22 - End")
	}
}

func otherOnlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		// t := time.Now()
		log.Printf("other group v2 - Satrt")
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("other group v2 - End")
	}
}

func main() {
	r := gee.New()
    r.Use(gee.Logger())

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV21(), onlyForV22(), otherOnlyForV2()) // v2 group middleware
    {
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
