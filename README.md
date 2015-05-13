golang-rest-api-server-example
==============================

[![Build Status](https://drone.io/github.com/sogko/golang-rest-api-server-example/status.png)](https://drone.io/github.com/sogko/golang-rest-api-server-example/latest)

A complete example of a REST-ful API server in Go

## Features
- Simple, flexible and testable architecture
-- Light-weight server component
-- Easy to replace components with your own library of choice (router, database driver, etc)
-- Guaranteed thread-safety for each request: uses `gorilla/context` for per-request context
- Not a framework
-- More like a project to quickly kick-start your own REST API server, customized to your own needs.
- Each REST resource is a separate package
-- Modular approach
-- Separation between `model`, `controller` and `data` layers
-- Clear abstraction from `server` package 
- Highly-testable code base
-- Unit-tested `server`; 100% code coverage
-- Easily test REST resources routes
-- Parallelizable test suite


## Quick start
```bash
# get go-package and put it in your go-workspace
go get github.com/sogko/golang-rest-api-server-example

# go to package root folder
cd $GOPATH/src/github.com/sogko/golang-rest-api-server-example

# install dependencies
go get

# run Server
go run main.go

# http://localhost:3001
```
-----

## Dependencies
- Golang v1.4+
- MongoDB
- External Go packages dependencies

```bash
# production
go get github.com/codegangsta/negroni   # HTTP server library
go get github.com/gorilla/mux           # HTTP router
go get github.com/gorilla/context       # Per-request context registry utility
go get github.com/unrolled/render       # JSON response renderer
go get gopkg.in/mgo.v2                  # Golang MongoDB driver

# development / test
go get github.com/onsi/ginkgo           # Golang BDD test framework, complements `go test`
go get github.com/onsi/gomega           # Ginkgo's preferred matcher library
go get github.com/modocache/gory        # `factory_girl` for Go
```

----

## Test
Install all Go package dependencies and run either one of the following command

```bash
go test

# or

$GOPATH/bin/ginkgo -r -p -node=4

# "-r" watches recursively (including test suites for sub-packages)
# "-p -nodes=4" parallelize test execution with 4 worker nodes
```

## TDD
```bash
$GOPATH/bin/ginkgo watch -r -p -nodes=4

# "-r" watches recursively (including test suites for sub-packages)
# "-p -nodes=4" parallelize test execution with 4 worker nodes
```

## Code coverage
To generate coverage profile

```bash
cd $GOPATH/src/github.com/sogko/golang-rest-api-server-example/server
$GOPATH/bin/ginkgo -cover
```

To view coverage

```bash
go tool cover -html=$GOPATH/src/github.com/sogko/golang-rest-api-server-example/server/server.coverprofile
```

## Sessions Management
To generate key pair for signing JWT claims
```bash
$ openssl genrsa -out demo.rsa 1024 # the 1024 is the size of the key we are generating
$ openssl rsa -in demo.rsa -pubout > demo.rsa.pub 
```

## TODO
* [x] API versioning using Accept header, for e.g: `Accept=application/json;version=1.0,*/*`
* [x] User and roles management
* [x] Session management using JWT tokens
* [x] Activity-based access control (ABAC)
* [-] Refactor using Dependency Inversion
* [ ] Task scheduler
* [ ] Load test using vegeta
* [ ] i18n (internationalisation)
* [ ] Implement another router library for kicks
* [ ] Consolidate util libraries and publish as separate package
* [x] Abstract away negroni middlewares to a generic http.HandlerFunc