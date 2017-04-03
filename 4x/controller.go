package main

import (
	"fmt"

	httpd "github.com/dharmadog888/opango/httpd"
)

/*
The X4Controller is an example Controller created
from the httpd package.
*/
type X4Controller struct {
	httpd.Controller
}

/*
NewX4Controller creates a new controller that will be created
by the app and assigned to the '/x4' Domain during the server start up and
assed to the router as (one of) the Domain Controllers for the app. The main
tast is to create the Routing by constructing all of the EndpointMaps that
connect an enpoint (URI), a method, and a processor .
*/
func NewX4Controller() *X4Controller {
	x4 := X4Controller{}

	// set up request routing for each of the
	// test cases for content type
	x4.SetupRouting([]httpd.EndpointMap{
		// JSON Payload
		*httpd.NewEndpointMap("/x4/json", httpd.POST, x4.handleJSON),
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
handleJSON demonstrates the JSON Passing. A little difficult from a single lightweight HTML page
so need to test via curl, python etc.
*/
func (x4 *X4Controller) handleJSON(req httpd.RestAPI, params httpd.ParamMap) (resp httpd.ApiResponse) {

	resp.Code = 200
	resp.Message = "Ok"
	resp.Body = fmt.Sprintf("%T - %v", req.JSON, req.JSON)

	return resp
}

/*
handleForm eturns a dump of the form data that came in from the request
*/
func (x4 *X4Controller) handleForm(req httpd.RestAPI, params httpd.ParamMap) (resp httpd.ApiResponse) {

	resp.Code = 200
	resp.Message = "Ok"
	resp.Body = fmt.Sprintf("%v", req.Form)

	return resp
}

/*
handleFile uploads a File from remote (client) computer (via file tag in HTML form)
*/
func (x4 *X4Controller) handleFile(req httpd.RestAPI, params httpd.ParamMap) (resp httpd.ApiResponse) {

	resp.Code = 200
	resp.Message = "Ok"
	resp.Body = "Files not yet implemented!"

	return resp
}

/*
handleQuerystring returns a dump of Querystring data sent with URL (standard k=v[&k=v...] format.
*/
func (x4 *X4Controller) handleQuerystring(req httpd.RestAPI, params httpd.ParamMap) (resp httpd.ApiResponse) {

	resp.Code = 200
	resp.Message = "Ok"
	resp.Body = fmt.Sprintf("QS %v", req.Querystring)

	return resp
}
