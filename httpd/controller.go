package httpd

/*
A Controller maps a REST endpoint to an
endpoint handler (processor)

handlers take the form

  handlerXXX(req httpd.RestAPI, params httpd.ParamMap) httpd.APIResponse

*/
import (
	"fmt"

	"regexp"
	"strings"
)

var (
	parmPtrn = regexp.MustCompile(`^\{(.+)\}$`)
)

const (
	pathSep = "/"
)

// ParamMap URI Path Parameter Map
type ParamMap map[string]string

// RestEndpoint function signature for API handlers
type RestEndpoint func(req RestAPI, params ParamMap) (resp APIResponse)

// EndpointSpec method and path for endpoint
type EndpointSpec struct {
	Method string
	Path   []string
}

// EndpointMap maps endpoint URI to handler
type EndpointMap struct {
	Endpoint  EndpointSpec
	Processor RestEndpoint
}

/*
NewEndpointMap is a convenience constructor for mapping endpoints
*/
func NewEndpointMap(path string, method string, proc RestEndpoint) *EndpointMap {
	em := EndpointMap{Endpoint: EndpointSpec{Method: method, Path: strings.Split(path, "/")}, Processor: proc}

	return &em
}

// IController The Controller interface spec
type IController interface {
	Route(params RestAPI) APIResponse
	SetupRouting(endpoints []EndpointMap) error
	Routes() []EndpointMap
}

/*
A Controller maps one or more REST endpoints to their assigned
endpoint handlers (processors)

handlers take the form

  handlerXXX(req httpd.RestAPI, params httpd.ParamMap) httpd.APIResponse

*/
type Controller struct {

	// the map of endpoints to methods
	// common to all controllers
	emp []EndpointMap
}

// Route direct request to correct processor
func (ctlr *Controller) Route(apiParams RestAPI) (resp APIResponse) {
	fmt.Println(apiParams)
	// test for match
	if params, proc := ctlr.URIMatch(apiParams.URI, apiParams.Method); proc != nil {
		// found one, call it...
		resp = proc(apiParams, params)

	} else {
		// not found...
		resp = APIResponse{}
		resp.Code = 404

		resp.Message = fmt.Sprintf("Path not found %s", apiParams.URL)
		resp.Body = "{\"Error\": \"No endpoint found\"}"
	}

	// return results
	return resp
}

// SetupRouting isntall controllers enpoint route mappings
func (ctlr *Controller) SetupRouting(endpoints []EndpointMap) (err error) {
	// save to internal map
	ctlr.emp = endpoints

	return err
}

// Routes reflect routing maps
func (ctlr Controller) Routes() []EndpointMap {
	// return mappings
	return ctlr.emp
}

// URIMatch tries to identify a processor given an array of URI nodes
func (ctlr Controller) URIMatch(uri []string, method string) (params map[string]string, processor RestEndpoint) {

	// check all endpoints for matching
	for tst := range ctlr.emp {
		var match bool

		ep := ctlr.emp[tst].Endpoint

		// method testing first...
		if !strings.EqualFold(method, ep.Method) {
			continue
		}
		// create a fresh param map
		pps := make(map[string]string)

		// fetch the stored path segments
		tendpt := ctlr.emp[tst].Endpoint

		// size match first
		if match = (len(tendpt.Path) == len(uri)); match {
			// segment by segment testing
			for seg := range uri {
				// check for param delimiters in endpts
				ptst := parmPtrn.FindStringSubmatch(tendpt.Path[seg])
				if len(ptst) > 0 {
					// found, create path param (implied match)
					pps[ptst[len(ptst)-1]] = uri[seg]
				} else {
					// not found, test path walk
					if !strings.EqualFold(uri[seg], tendpt.Path[seg]) {
						// not a match
						match = false
						break
					}
				}
			}
		}

		// check results
		if match {
			// winner winner chicken dinner!
			params = pps
			processor = ctlr.emp[tst].Processor
			break
		}
	}

	// and the results are....
	return params, processor
}

// Response conveniece response enerator
func (ctlr Controller) Response(code int, message string, body string) (resp APIResponse) {
	resp.Code = code
	resp.Message = message
	resp.Body = body

	return resp
}

// APIError conveniece error enerator
func (ctlr Controller) APIError(code int, err error) APIResponse {
	return ctlr.Response(code, err.Error(), "")
}
