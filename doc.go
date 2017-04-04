/*
Opango is a MVC toolkit that extends the standard Go packages with some
nifty things intended to accelerate microservice development in Go. The
Opango flagship is the pango/httpd package which wraps the standard Go
HTTP(s) servers (net/http package) with a RESTful router / controller
architecture. See the opango/4x package to see how to how this all goes
together. The typical server can usually be done in under 150 lines and
you go from there straight into a Controller that demonstrates many of the
request parsing features focusing on the various Content-Types supported
(JSON, Form Data and File Upload mime types)
*/
package opango
