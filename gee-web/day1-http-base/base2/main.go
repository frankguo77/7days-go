package  main


import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct{}


func (eg *Engine) ServeHTTP(w http.ResponseWritter, req *http.Request) {
	switch req.URL.Path {
        case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", requ.URL.Path)
	case
	}
}
