package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	x4 "github.com/dharmadog888/opango/4x/httpd4x"
	httpd "github.com/dharmadog888/opango/httpd"
)

const (
	// defaultConfig points to config file containing
	// server ip:port and location of html content
	defaulfConfig = "./config/x4.json"
)

var (
	// serverConf Server configuration poulated from config file
	config serverConf
)

// serverConf configuration matters
type serverConf struct {
	Addr        string
	ContentRoot string
}

// Read and parse JSON config file
func initServerConf(path string) (err error) {
	// read...
	data, e := ioutil.ReadFile(path)
	fmt.Printf("%s\n", data)
	if e == nil {
		// ...parse
		json.Unmarshal(data, &config)
		fmt.Printf("%v", config)
	}
	return err
}

// Set up routing with a handler wrapper
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

// Usage for command line
func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

//---------------------
// Web server PEP
func main() {
	// command line processing
	flag.Usage = Usage
	var confLoc = flag.String("config", defaulfConfig, "Location of Server Config file (./config/4x.json)")
	flag.Parse()

	// set up the config info
	if err := initServerConf(*confLoc); err != nil {
		log.Fatalf("Fatal Config Error: %s", err)
	}

	// build the web router with a single controller and an html content root
	router := httpd.NewRouter(httpd.RouteMap{"/x4": x4.NewX4Controller()}, config.ContentRoot)

	// configure the server
	s := &http.Server{
		Addr:           config.Addr,
		Handler:        makeHandler(router.RouteRequest),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// engage!
	log.Printf("-- Starting Web Service on %s", config.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("!! Web Server Crashed:\n--> %s", err)
	}
}
