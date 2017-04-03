package httpd

/*
   A Controller maps a REST endpoint to an
   endpoint processor
*/
import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	parmPtrn = regexp.MustCompile(`^\{(.+)\}$`)
)

const (
	pathSep = "/"
)

// URI Path Parameter Map
type ParamMap map[string]string

// Endpoint Processor - entrance to the Middle Kingdom
type RestEndpoint func(req RestAPI, params ParamMap) (resp ApiResponse)

// Endpoint Specification
type EndpointSpec struct {
	Method string
	Path   []string
}

// Routing Specification
type EndpointMap struct {
	Endpoint  EndpointSpec
	Processor RestEndpoint
}

func NewEndpointMap(path string, method string, proc RestEndpoint) *EndpointMap {
	em := EndpointMap{Endpoint: EndpointSpec{Method: method, Path: strings.Split(path, "/")}, Processor: proc}

	return &em
}

// Controller Spec
type IController interface {
	Route(params RestAPI) ApiResponse
	SetupRouting(endpoints []EndpointMap) error
	Routes() []EndpointMap
}

// Base implementation type
type Controller struct {

	// the map of endpoints to methods
	// common to all controllers
	emp []EndpointMap
}

// direct request to correct processor
func (ctlr *Controller) Route(apiParams RestAPI) (resp ApiResponse) {

	// test for match
	if params, proc := ctlr.UriMatch(apiParams.Uri, apiParams.Method); proc != nil {
		// found one, call it...
		resp = proc(apiParams, params)

	} else {
		// not found...
		resp = ApiResponse{}
		resp.Code = 404

		resp.Message = fmt.Sprint("Path not found	(%s)", apiParams.Uri)
		resp.Body = "{\"Error\": \"No endpoint found\"}"

		log.Printf("?? Routing not found %v\n", apiParams)
	}

	// return results
	return resp
}

// set up the controllers enpoint route mappings
func (ctlr *Controller) SetupRouting(endpoints []EndpointMap) (err error) {
	// save to internal map
	ctlr.emp = endpoints

	return err
}

// reflect routing maps
func (ctlr Controller) Routes() (emp []EndpointMap) {
	// return mappings
	return emp
}

// the way
func (ctlr Controller) UriMatch(uri []string, method string) (params map[string]string, processor RestEndpoint) {

	// optimistic sets
	match := true

	// check all endpoints for matching
	for tst := range ctlr.emp {
		ep := ctlr.emp[tst].Endpoint

		// method testing first...
		if !strings.EqualFold(method, ep.Method) {
			match = false
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

// Response Generator
func (ctlr Controller) Response(code int, message string, body string) (resp ApiResponse) {
	resp.Code = code
	resp.Message = message
	resp.Body = body

	return resp
}

func (ctlr Controller) ApiError(code int, err error) ApiResponse {
	return ctlr.Response(code, err.Error(), "")
}
