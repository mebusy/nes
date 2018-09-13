package cheater

import (
	_ "encoding/binary"
	"fmt"
	"io"
	_ "io/ioutil"
	"log"
	"nes/nes"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

var ramMap sync.Map

func readMem(mem nes.Memory, startAddr, nbyte int) uint32 {
	var val uint32 = 0
	for i := 0; i < nbyte; i++ {
		val += uint32(mem.Read(uint16(startAddr+i))) << uint32(8*i)
	}
	return val
}
func clearMap(mem nes.Memory) {
	delete := func(key interface{}, value interface{}) bool {
		ramMap.Delete(key)
		return true
	}
	ramMap.Range(delete)

	for addr := 0; addr < 0x8000; addr++ {
		// 2K RAM , and 8K SRAM
		for nbyte := 1; nbyte <= 4; nbyte++ {
			// 1 2 3 4
			if addr+nbyte-1 >= 2048 && addr+nbyte-1 < 0x6000 {
				continue
			}
			key := fmt.Sprintf("%d:%d", addr, nbyte)
			val := readMem(mem, addr, nbyte)
			store2map(key, val)
		}
	}

	// test
	/*
		for i := 0; i < 4; i++ {
			fmt.Printf("%02X ", readMem(mem, 0x6BF5+i, 1))
		}
		key := fmt.Sprintf("%d:%d", 0x6BF5, 4)
		val, ok := ramMap.Load(key)
		fmt.Printf(" %X, %v", val, ok)
		fmt.Println()
		//*/
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

func print10matchedAddress(w http.ResponseWriter, mem nes.Memory) {
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
			// TODO out range
			io.WriteString(w, fmt.Sprintf("%02x ", readMem(mem, addr+i, 1)))
		}
		io.WriteString(w, "\n")

		return true
	}
	ramMap.Range(dump)

}

func search(mem nes.Memory, w http.ResponseWriter, req *http.Request) {
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

			curVal := readMem(mem, addr, nByte)
			if comparator(val, curVal) {
				store2map(key, curVal)
			} else {
				ramMap.Delete(key)
			}

			return true
		}
		ramMap.Range(valueCompare)

		return
	}

	// compare stored value and current value
	vals, ok = req.URL.Query()["diff"]
	if ok && len(vals[0]) >= 1 {
		op := vals[0]
		var comparator func(val1, val2 interface{}) bool

		switch op {
		case "inc":
			comparator = func(val1, val2 interface{}) bool {
				v1 := interface2uint32(val1)
				v2 := interface2uint32(val2)
				return v1 > v2
			}
		case "dec":
			comparator = func(val1, val2 interface{}) bool {
				v1 := interface2uint32(val1)
				v2 := interface2uint32(val2)
				return v1 < v2
			}
		case "eq":
			comparator = func(val1, val2 interface{}) bool {
				v1 := interface2uint32(val1)
				v2 := interface2uint32(val2)
				return v1 == v2
			}
		case "neq":
			comparator = func(val1, val2 interface{}) bool {
				v1 := interface2uint32(val1)
				v2 := interface2uint32(val2)
				return v1 != v2
			}
		default:
			io.WriteString(w, "incorrect diff argument \n")
			return
		}

		valueCompare := func(_key interface{}, _ interface{}) bool {
			addr := 0
			nByte := 0
			key := _key.(string)
			fmt.Sscanf(key, "%d:%d", &addr, &nByte)

			_val, ok2 := ramMap.Load(key)
			storedVal := _val.(uint32)
			curVal := readMem(mem, addr, nByte)

			if ok2 && comparator(storedVal, curVal) {
				store2map(key, curVal)
			} else {
				ramMap.Delete(key)
			}

			return true
		}
		ramMap.Range(valueCompare)

		return
	}

	io.WriteString(w, "incorrect query \n")

}

func showMem(mem nes.Memory, w http.ResponseWriter, req *http.Request) {
	vals, ok := req.URL.Query()["addr"]
	if ok && len(vals[0]) >= 1 {
		if addr, err := strconv.ParseUint(vals[0], 16, 16); err == nil {
			io.WriteString(w, fmt.Sprintf("%04X :  ", addr))
			for i := 0; i < 4; i++ {
				// TODO out range
				io.WriteString(w, fmt.Sprintf("%02x ", readMem(mem, int(addr)+i, 1)))
			}
			io.WriteString(w, "\n")
			return
		}
	}

	io.WriteString(w, "incorrect query \n")
}

func setMem(mem nes.Memory, w http.ResponseWriter, req *http.Request) {
	addrs, ok1 := req.URL.Query()["addr"]
	vals, ok2 := req.URL.Query()["val"]
	nBytes, ok3 := req.URL.Query()["n"]

	if ok1 && ok2 && ok3 {
		if addr, err := strconv.ParseUint(addrs[0], 16, 16); err == nil {
			if nByte, err := strconv.ParseUint(nBytes[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(vals[0], 10, int(nByte*8)); err == nil {
					for i := 0; i < int(nByte); i++ {
						b := val & 0xFF
						mem.Write(uint16(int(addr)+i), byte(b))
						val >>= 8
					}
				}
			}
		}
		io.WriteString(w, "cheated \n")
		return
	}
	io.WriteString(w, "incorrect query \n")
}

func StartWebServer(console *nes.Console) {
	// Hello world, the web server

	go func() {
		helloHandler := func(w http.ResponseWriter, req *http.Request) {
			//fmt.Printf("receive requets: %s - %s\n", req.Method, req.URL.Path)
			//ram := console.RAM
			if req.URL.Path == "/clear" {
				clearMap(console.CPU.Memory)
			} else if req.URL.Path == "/search" {
				search(console.CPU.Memory, w, req)
			} else if req.URL.Path == "/show" {
				showMem(console.CPU.Memory, w, req)
			} else if req.URL.Path == "/set" {
				setMem(console.CPU.Memory, w, req)
			} else {
				io.WriteString(w, "Example\n")
				io.WriteString(w, "/clear\n")
				io.WriteString(w, "/search?val=12\n")
				io.WriteString(w, "/search?diff=(inc|dec|eq|neq)\n")
				io.WriteString(w, "/show?addr=6076\n")
				io.WriteString(w, "/set?addr=6076&val=10000&n=3\n")
			}

			//
			io.WriteString(w, "matched result: ")
			io.WriteString(w, fmt.Sprintf("%d\n\n", mapCount()))
			print10matchedAddress(w, console.CPU.Memory)
		}

		http.HandleFunc("/", helloHandler)
		log.Fatal(http.ListenAndServe(":7777", nil))
	}()
}
