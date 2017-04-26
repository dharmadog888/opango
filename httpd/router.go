package httpd

/*
	RESTful Routing and Content Services

	The web router construction takes a routing map of
	URIs (endpoints) and Controllers along with a web root folder
	name and sets up for the actually business of serving files
	and data

	TODO: Implement list of index pages to check against...
*/

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	defaultIndexPage = "index.html"
	ioBuffSize       = 2 * 1024 * 1024 // 2M
)

/*
POST, GET, PUT, DELETE httpd method names
*/
const (
	POST = "POST"
	GET  = "GET"
	PUT  = "PUT"
	HEAD = "HEAD"
)

// RestAPI struct sent to controller that describes request
type RestAPI struct {
	Method      string
	URL         string
	URI         []string
	JSON        json.Decoder
	Form        url.Values
	Header      http.Header
	MPForm      multipart.Form
	Querystring url.Values
	Body        string
}

// HTTP Results
type httpResponse struct {
	Code    int
	Message string
	Cookie  http.Cookie
}

// APIResponse container for results
type APIResponse struct {
	httpResponse
	Body string
}

// RouteMap map to assign URI to Controller
type RouteMap map[string]IController

/*
Router mananges requests for services
on behalf of an App Server
*/
type Router struct {
	routes     RouteMap
	webroot    string
	indexPages []string
}

// NewRouter constructor for  new Router
func NewRouter(rmap RouteMap, root string) Router {
	ndxpgs := make([]string, 1)
	ndxpgs[0] = defaultIndexPage

	rtr := Router{routes: rmap, webroot: root, indexPages: ndxpgs}

	return rtr
}

// SetRoute add/update uri to route mappings
func (r Router) SetRoute(uri string, ctlr Controller) {
	r.routes[uri] = &ctlr
}

// ClearRoute remove uri mapping
func (r Router) ClearRoute(uri string) {
	delete(r.routes, uri)
}

//  routes are defined by the first node(s) of the uri
func (r Router) findRoute(uri string) (ctlr IController, found string) {
	// concation buffer to avoid
	// string churning
	var tst bytes.Buffer

	// path nodes
	parts := strings.Split(strings.Split(uri, "?")[0], "/")
	if !strings.HasPrefix(uri, "/") {
		// no relative paths in API
		return nil, ""
	}

	// absolute pathing from root
	parts[0] = "/"

	// hunt for first match
	for p := range parts {
		if p > 1 {
			// need path delimiter?
			tst.WriteString("/")
		}

		// add node
		tst.WriteString(parts[p])

		// test
		var good bool
		if ctlr, good = r.routes[tst.String()]; good {
			found = tst.String()

			// break on first match
			break
		}
	}

	// return results
	return ctlr, found
}

// web server 101
func (r Router) serveFile(path string, resp http.ResponseWriter) error {
	var err error
	var body []byte

	// file exists?
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return errors.New("File Not Found")
	}

	// read and write
	if body, err = ioutil.ReadFile(path); err != nil {
		return err
	}

	log.Printf("++ <-- 200 Ok: %s\n", path)
	resp.Write(body)

	return err
}

// RouteRequest Start point for routing from server ** //
func (r Router) RouteRequest(resp http.ResponseWriter, req *http.Request) {

	reqURL := fmt.Sprintf("%s", req.URL)

	// hacking filters
	if strings.Contains(reqURL, "..") {
		log.Printf("!! <-- 402: Forbidden (%s)", reqURL)
		http.Error(resp, "Forbidden,", 402)
		return
	}

	ctrl, found := r.findRoute(reqURL)

	// do we have routing
	if found != "" {
		// sanity check.... (aka: shouldn't happen but...)
		if ctrl == nil {
			log.Println("!! <-- 500: Invalid controller ", found)
			http.Error(resp, "Invalid Controller", 500)
			return
		}

		// strip out found nodes
		uri := strings.Split(reqURL[len(found):], "/")

		// and delegate to found controller
		api := RestAPI{URL: reqURL, URI: uri, Method: req.Method, Header: req.Header, Querystring: req.URL.Query()}

		// some investigation is required to
		// figure out what type of data we got
		var hdr []string
		if hdr = req.Header["Content-Type"]; hdr == nil {
			hdr = []string{"text/plain"}
		}

		switch ctype := hdr[0]; ctype {
		case "application/json":
			// JSON
			if req.Body != nil {
				api.JSON = *json.NewDecoder(req.Body)
			}

		case "multipart/form":
			// File Upload
			if err := req.ParseMultipartForm(ioBuffSize); err == nil {
				api.MPForm = *req.MultipartForm
			}

		case "application/x-www-form-urlencoded":
			// Form data
			req.ParseForm()
			api.Form = req.Form

		default:
			// parse body as text
			if req.Body != nil {
				buff := make([]byte, ioBuffSize)
				for {
					var err error
					var cnt int
					if cnt, err = req.Body.Read(buff); cnt > 0 {
						api.Body = fmt.Sprintf("%s%s", api.Body, buff[:cnt])
					}
					if err == io.EOF {
						break
					} else if err != nil {
						log.Println("!! <-- IO Error reading request body", found)
						http.Error(resp, "IO Error", 500)
						return
					}
				}
			}
		}
		var apiRsp APIResponse
		if apiRsp = ctrl.Route(api); &apiRsp.Cookie != nil {
			http.SetCookie(resp, &apiRsp.Cookie)
		}

		// what say the controler?
		switch {
		case apiRsp.Code > 399:
			log.Printf("-- <-- %d Error: %s", apiRsp.Code, apiRsp.Message)
			http.Error(resp, apiRsp.Message, apiRsp.Code)
		case apiRsp.Code > 299:
			log.Printf("~~ <-- %d Redirect: %s", apiRsp.Code, apiRsp.Message)
			http.Redirect(resp, req, apiRsp.Message, apiRsp.Code)
		case apiRsp.Code == 200:
			log.Printf("~~ <-- %d Ok: %s", apiRsp.Code, apiRsp.Message)
			resp.Write([]byte(apiRsp.Body))
		}

	} else {
		// no, treat as content chained to web root
		var path string

		// default index page for root request
		if reqURL == "/" {
			reqURL = fmt.Sprintf("/%s", r.indexPages[0])
		}

		// keep absolute path "chrooted"
		if strings.HasPrefix(reqURL, "/") {
			path = r.webroot
		} else {
			// relative path
			pathnodes := strings.Split(req.Referer(), "/")
			path = strings.Join(pathnodes[0:len(pathnodes)-1], "/")
		}

		// set real path to file
		path += reqURL

		// and attempt to serve it up
		if err := r.serveFile(path, resp); err != nil {
			log.Printf("-- <-- 404 - File not found: %s", path)
			http.NotFound(resp, req)
		}
	}
}
