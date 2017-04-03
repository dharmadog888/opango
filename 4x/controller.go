package main

import (
	"fmt"

	httpd "bitbucket.org/dharmadog888/opango/httpd"
)

/*
  Example Controller created
  witn the httpd package.
*/
type X4Controller struct {
	httpd.Controller
}

/*
  X4 Controller Constructor.

  This controller will be created by the app and assigned to the
  '/x4' Domain during the server start up and passed to the router
  as (one of) the Domain Controllers for the app. The main
  tast is to create the Routing by constructing all of the
  EndpointMaps that connect an enpoint (URI), a method, and
  a processor .

  TODO: Currently a separate map for each endpoint/method
        combo is required. Future phase will allow for multiple
		methods per endpoint to make this a bit easier in those
		situations.

*/
func NewX4Controller() *X4Controller {
	x4 := X4Controller{}

	// set up request routing for each of the
	// test cases for content type
	x4.SetupRouting([]httpd.EndpointMap{
		// JSON Payload
		*httpd.NewEndpointMap("/x4/json", httpd.POST, x4.handleJson),
		// HTML Form Data
		*httpd.NewEndpointMap("/x4/form", httpd.POST, x4.handleForm),
		// Querystring
		*httpd.NewEndpointMap("/x4/qs", httpd.GET, x4.handleQuerystring),
		// File Upload
		*httpd.NewEndpointMap("/x4/file", httpd.POST, x4.handleFile),
	})

	// return completed controller
	return &x4
}

/*
  Demonstrated the JSON Passing. A little difficult from a single lightweight HTML page
  so need to test via curl, python etc.
*/
func (x4 *X4Controller) handleJson(req httpd.RestAPI, params httpd.ParamMap) (resp httpd.ApiResponse) {

	resp.Code = 200
	resp.Message = "Ok"
	resp.Body = fmt.Sprintf("%T - %V", req.JSON, req.JSON)

	return resp
}

/*
  Returns a dump of the form data that came in from the request
*/
func (x4 *X4Controller) handleForm(req httpd.RestAPI, params httpd.ParamMap) (resp httpd.ApiResponse) {

	resp.Code = 200
	resp.Message = "Ok"
	resp.Body = fmt.Sprintf("%v", req.Form)

	return resp
}

/*
	Upload a File from remote (client) computer (via file tag in HTML form)
*/
func (x4 *X4Controller) handleFile(req httpd.RestAPI, params httpd.ParamMap) (resp httpd.ApiResponse) {

	resp.Code = 200
	resp.Message = "Ok"
	resp.Body = "Files not yet implemented!"

	return resp
}

/*
  Returns a dump of Querystring data sent with URL (standard k=v[&k=v...] format.
*/
func (x4 *X4Controller) handleQuerystring(req httpd.RestAPI, params httpd.ParamMap) (resp httpd.ApiResponse) {

	resp.Code = 200
	resp.Message = "Ok"
	resp.Body = fmt.Sprintf("QS %v", req.Querystring)

	return resp
}
