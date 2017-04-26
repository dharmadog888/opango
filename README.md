# Opango

Introduction
------------

Opango is the open source branch of the Pangolin project. We use this 
project to expose some nifty things that think might be useful for the Go 
gang in general. The prime motivators for the opening are contained in the 
httpd package and consists of a REST based web routing / controller system 
that plugs into the golang net/http package to implement RESTful APIs.

Compatibility
-------------
 *Opango* is intended to work with the standard Golang net/http library and 
 has been tested with Go (major) versions 1.7 and 1.8

Installation
------------

To get opango, run

```
 $> go get github.com/dharamdog888/opango/...
 (for entire project including 4x )

 $> go get github.com/dharmadog888/opango/httpd
 (for just the Web Routing tools)
```

License
-------

The opnago package is licensed under the Apache License 2.0. Please see the LICENSE file for details.

Example
-------
Tour the 4x package for sample server, there is a working server and
controller that demonstrate getting a RESTful service up and running.


The Buck Stops Here
-------------------

 * David Connell <dharmadog888@gmail.com>
