package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"nes/ui"
    "flag"
    "runtime/pprof"

	"io"
	"net/http"
    "nes/cheater"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var romPath = flag.String("p", "", "nes rom path")

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        log.Println( "output profile to " , *cpuprofile )
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal("could not create CPU profile: ", err)
        }
        if err := pprof.StartCPUProfile(f); err != nil {
            log.Fatal("could not start CPU profile: ", err)
        }
        defer func() { 
            pprof.StopCPUProfile() 
            log.Println( "stop" )
        }()
    }

    cheater.StartWebServer()

	log.SetFlags(0)
	// paths := getPaths()
    paths := []string{ *romPath}


	if len(paths) == 0 {
		log.Fatalln("no rom files specified or found")
	}
	ui.Run(paths)
}

func getPaths() []string {
	var arg string
	args := os.Args[1:]
	if len(args) == 1 {
		arg = args[0]
	} else {
		arg, _ = os.Getwd()
	}
	info, err := os.Stat(arg)
	if err != nil {
		return nil
	}
	if info.IsDir() {
		infos, err := ioutil.ReadDir(arg)
		if err != nil {
			return nil
		}
		var result []string
		for _, info := range infos {
			name := info.Name()
			if !strings.HasSuffix(name, ".nes") {
				continue
			}
			result = append(result, path.Join(arg, name))
		}
		return result
	} else {
		return []string{arg}
	}
}

func startWebServer() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":7777", nil))    
}
