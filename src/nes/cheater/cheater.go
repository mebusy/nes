package cheater

import (
	"encoding/binary"
	"fmt"
	"io"
	_ "io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

var ramMap sync.Map

func clearMap(ram []byte) {
	delete := func(key interface{}, value interface{}) bool {
		ramMap.Delete(key)
		return true
	}
	ramMap.Range(delete)

	for addr, v := range ram {
		//fmt.Println(addr, v)
		key := fmt.Sprintf("%d:%d", addr, 1)
		store2map(key, v)

		if addr <= len(ram)-2 {
			key = fmt.Sprintf("%d:%d", addr, 2)
			store2map(key, binary.LittleEndian.Uint16(ram[addr:]))
		}
		if addr <= len(ram)-4 {
			key = fmt.Sprintf("%d:%d", addr, 4)
			store2map(key, binary.LittleEndian.Uint32(ram[addr:]))
		}

	}

}

func mapCount() int {
	cnt := 0
	count := func(key interface{}, value interface{}) bool {
		cnt++
		return true
	}
	ramMap.Range(count)
	return cnt
}

func interface2uint32(val interface{}) uint32 {
	switch val.(type) {
	case int:
		return uint32(val.(int))
	case uint8:
		return uint32(val.(uint8))
	case uint16:
		return uint32(val.(uint16))
	case uint32:
		return val.(uint32)
	default:
		fmt.Println(reflect.ValueOf(val).Kind())
		panic("不知道什么整数类型")
	}

}

func store2map(key string, val interface{}) {
	ramMap.Store(key, interface2uint32(val))
}

func debugInfo(w http.ResponseWriter, ram []byte) {
	//io.WriteString(w, fmt.Sprintf("val of addr 0: %d", ram[0]))
}

func print10matchedAddress(w http.ResponseWriter, ram []byte) {
	cnt := 0
	dump := func(_key interface{}, _value interface{}) bool {
		addr := 0
		nByte := 0
		key := _key.(string)
		fmt.Sscanf(key, "%d:%d", &addr, &nByte)

		cnt++
		if cnt > 10 {
			return false
		}

		io.WriteString(w, fmt.Sprintf("%04X : %12d , ", addr, _value.(uint32)))
		for i := 0; i < nByte; i++ {
			io.WriteString(w, fmt.Sprintf("%02x ", ram[addr+i]))
		}
		io.WriteString(w, "\n")

		return true
	}
	ramMap.Range(dump)

}

func search(ram []byte, w http.ResponseWriter, req *http.Request) {
	vals, ok := req.URL.Query()["val"]

	if ok && len(vals[0]) >= 1 {
		// Query()["val"] will return an array of items,
		// we only want the single item.
		val, _ := strconv.Atoi(vals[0])

		// only compare current value  with input value , no stored value
		equal := func(val1, val2 interface{}) bool {
			v1 := interface2uint32(val1)
			v2 := interface2uint32(val2)
			return v1 == v2
		}
		comparator := equal

		valueCompare := func(_key interface{}, _ interface{}) bool {
			addr := 0
			nByte := 0
			key := _key.(string)
			fmt.Sscanf(key, "%d:%d", &addr, &nByte)

			if nByte == 1 && comparator(val, ram[addr]) {
				// store new value
				store2map(key, ram[addr])

			} else if nByte == 2 && comparator(val, binary.LittleEndian.Uint16(ram[addr:])) {
				store2map(key, binary.LittleEndian.Uint16(ram[addr:]))
			} else if nByte == 4 && comparator(val, binary.LittleEndian.Uint32(ram[addr:])) {
				store2map(key, binary.LittleEndian.Uint32(ram[addr:]))
			} else {
				ramMap.Delete(key)
			}

			return true
		}
		ramMap.Range(valueCompare)

		io.WriteString(w, "matched result: ")
		io.WriteString(w, fmt.Sprintf("%d\n", mapCount()))
		debugInfo(w, ram)

		return
	}

	// compare stored value and current value
	vals, ok = req.URL.Query()["diff"]
	if ok && len(vals[0]) >= 1 {
		op := vals[0]
		_ = op

		lessThan := func(val1, val2 interface{}) bool {
			v1 := interface2uint32(val1)
			v2 := interface2uint32(val2)
			return v1 < v2
		}
		comparator := lessThan

		valueCompare := func(_key interface{}, _ interface{}) bool {
			addr := 0
			nByte := 0
			key := _key.(string)
			fmt.Sscanf(key, "%d:%d", &addr, &nByte)

			_val, ok2 := ramMap.Load(key)
			val := _val.(uint32)

			if nByte == 1 && ok2 && comparator(val, ram[addr]) {
				// store new value
				store2map(key, ram[addr])
			} else if nByte == 2 && comparator(val, binary.LittleEndian.Uint16(ram[addr:])) {
				store2map(key, binary.LittleEndian.Uint16(ram[addr:]))
			} else if nByte == 4 && comparator(val, binary.LittleEndian.Uint32(ram[addr:])) {
				store2map(key, binary.LittleEndian.Uint32(ram[addr:]))
			} else {
				ramMap.Delete(key)
			}

			return true
		}
		ramMap.Range(valueCompare)

		io.WriteString(w, "matched result: ")
		io.WriteString(w, fmt.Sprintf("%d", mapCount()))

		return
	}

	io.WriteString(w, "incorrect query \n")

}

func StartWebServer(ram []byte) {
	// Hello world, the web server

	go func() {
		helloHandler := func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/clear" {
				clearMap(ram)
				io.WriteString(w, "clear!\n")
				debugInfo(w, ram)
			} else if req.URL.Path == "/search" {
				search(ram, w, req)
			} else {
				io.WriteString(w, "Example\n")
				io.WriteString(w, "/clear\n")
				io.WriteString(w, "/search?val=12\n")
			}

			//
			print10matchedAddress(w, ram)
		}

		http.HandleFunc("/", helloHandler)
		log.Fatal(http.ListenAndServe(":7777", nil))
	}()
}
