/*
The opango/httpd package supplies tools for implementing RESTful API
services by enhancing the Go standard *net/http* packages. The Router is
ready to use as is and the Controller provides the basis for managing
endpoint processors.

The general concept is to create a RouteMap than assigns URI Domains, which
can be thought of as the first node(s) or "prefix" of a URI, to a
specic Controller. For example, a "/cipher" Domain may host the endpoints
"/cipher/encrypt" and "/cipher/decrypt" that are both implemented by a
CipherController. The RouteMap would associate "/cipher" with an instance of
a CipherController so all incoming requests that have a URI starting with
"/cipher" would be delegated to the CipherController for processing. The
CipherController uses and EndpointMap to further route the request by exact
matching of the full URI and the HTTP Method of the incoming request, to a
function on the Controller.

During the processing, the incoming request is held by the Router which uses
the Content-Type to extract any incoming payload and make it available to
the Controller in the format expected. The Controller will return a
APIResponse back to the router which marshals the results into the native
HTTP ResponseWriter.

The dharmadog888/opango/4x package includes a simple server and controller
to demonstrate how this all goes together.
*/
package httpd
