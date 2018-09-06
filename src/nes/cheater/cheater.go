package cheater

import (
	"io"
    "io/ioutil"
	"log"
	"net/http"
)

func StartWebServer() {
	// Hello world, the web server

    go func() {
        helloHandler := func(w http.ResponseWriter, req *http.Request) {
            io.WriteString(w, "Hello, world!\n"  )
            io.WriteString(w,  req.URL.Path )
        }

        http.HandleFunc("/", helloHandler)
        log.Fatal(http.ListenAndServe(":7777", nil))
    }() 
}

