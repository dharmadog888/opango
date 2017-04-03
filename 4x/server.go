/*
  4x - For Example web app...

  Demonstrates using the httpd package
  resources to rapidly build a working
  RESTful web server.

  Routing is managed by the httpd.Router and
  is built specifically to handle the task of
  implementing microservices via an HTTP REST
  protocol. Due to the nature of APIs in general
  and microserives specifally, this has been
  limited to the following transaction types
  as recognized via the HTTP Content-Type header:

    Form Data: application/x-www-form-urlencoded
	JSON Data: application/json
	File Uploads: multipart/form-data

  All other mime types are treated as text when
  sent to the controller and can be specifically
  managed at the application level.
*/
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

	httpd "bitbucket.org/dharmadog888/opango/httpd"
)

const (
	// server ip:port and location of html content
	DefaultConfig = "./config/x4.json"
)

var (
	// Server configuration
	serverConf ServerConf
)

// Server configuration matters
type ServerConf struct {
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
		json.Unmarshal(data, &serverConf)
		fmt.Printf("%v", serverConf)
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
	var confLoc = flag.String("config", DefaultConfig, "Location of Server Config file (./config/4x.json)")
	flag.Parse()

	// set up the config info
	if err := initServerConf(*confLoc); err != nil {
		log.Fatalf("Fatal Config Error: %s", err)
	}

	// build the web router with a single controller and an html content root
	router := httpd.NewRouter(httpd.RouteMap{"/x4": NewX4Controller()}, serverConf.ContentRoot)

	// configure the server
	s := &http.Server{
		Addr:           serverConf.Addr,
		Handler:        makeHandler(router.RouteRequest),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// engage!
	log.Printf("-- Starting Web Service on %s", serverConf.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("!! Web Server Crashed:\n--> %s", err)
	}
}
