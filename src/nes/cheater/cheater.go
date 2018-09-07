package cheater

import (
	"encoding/binary"
	"fmt"
	"io"
	_ "io/ioutil"
	"log"
	"net/http"
	"sync"
)

var ramMap sync.Map

func clearMap(ram []byte) {
	delete := func(key interface{}, value interface{}) bool {
		ramMap.Delete(key)
		return true
	}
	ramMap.Range(delete)

	for i, v := range ram {
		// log.Println(i, v)
		key := fmt.Sprintf("%d:%d", i, 1)
		ramMap.Store(key, v)

		if i < len(ram)-1 {
			key = fmt.Sprintf("%d:%d", i, 2)
			ramMap.Store(key, binary.LittleEndian.Uint16(ram[i:]))
		}
		if i < len(ram)-3 {
			key = fmt.Sprintf("%d:%d", i, 4)
			ramMap.Store(key, binary.LittleEndian.Uint32(ram[i:]))
		}

	}

}

func search(ram []byte, w http.ResponseWriter, req *http.Request) {
	vals, ok := req.URL.Query()["val"]

	if ok && len(vals[0]) >= 1 {
		// Query()["val"] will return an array of items,
		// we only want the single item.
		val := vals[0]

		// only compare current value  , and compared value , no stored value
		io.WriteString(w, "val: ")
		io.WriteString(w, val)
		return
	}

	io.WriteString(w, "incorrect query \n")

}

func StartWebServer(ram []byte) {
	// Hello world, the web server

	go func() {
		helloHandler := func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/clear" {
				io.WriteString(w, "clear!\n")
				clearMap(ram)
			} else if req.URL.Path == "/search" {
				search(ram, w, req)

			} else {
				io.WriteString(w, "Example\n")
				io.WriteString(w, "/clear\n")
				io.WriteString(w, "/search?val=12\n")
			}
		}

		http.HandleFunc("/", helloHandler)
		log.Fatal(http.ListenAndServe(":7777", nil))
	}()
}
