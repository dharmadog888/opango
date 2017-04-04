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
