[![slumber](http://i.imgur.com/RXDVdB0.png)](https://github.com/sogko/slumber)
[![Build Status](https://drone.io/github.com/sogko/slumber/status.png)](https://drone.io/github.com/sogko/slumber/latest)
[![Coverage Status](https://coveralls.io/repos/sogko/slumber/badge.svg?branch=master)](https://coveralls.io/r/sogko/slumber?branch=master)

A complete example of a REST-ful API server in written in Go (golang).

-----

## Features
- Simple, flexible and testable architecture
  - Light-weight server component
  - Easy to replace components with your own library of choice (router, database driver, etc)
  - Guaranteed thread-safety for each request: uses `gorilla/context` for per-request context
  - Uses dependency-inversion principle (DI) to reduce complexity and manage package dependencies.
- Not a framework
  - More like a project to quickly kick-start your own REST API server, customized to your own needs.
  - Easily extend the project with your own REST resources, in addition to the built-in `users` and `sessions` resources.
- Each REST resource is a separate package
  - Modular approach
  - Separation between `model`, `controller` and `data` layers
  - Clear abstraction from `server` package 
  - Take a look at built-in resources for examples: `users` and `sessions`
  - More example projects coming soon!
- Batteries come included
  - API versioning using using Accept header, for e.g: `Accept=application/json;version=1.0,*/*`
  - Default resources for `users` and `sessions`
  - Access control using activity-based access control (ABAC)
  - Authentication and session management using JWT token
  - Context middleware using `gorilla/context` for per-request context
  - JSON response rendering using `unrolled/render`; extensible to XML or other formats for response
  - MongoDB middleware for database; extensible for other database drivers
- Highly-testable code base
  - Unit-tested `server`; 100% code coverage
  - Easily test REST resources routes
  - Parallelizable test suite
  - Uses `ginkgo` for test framework; optional.

## Changelog
See [CHANGELOG.md](./CHANGELOG.md) for changes

## Quick start
To run an instance of a server example:

```bash
# get go-package and put it in your go-workspace
go get github.com/sogko/slumber

# go to package root folder
cd $GOPATH/src/github.com/sogko/slumber

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
go get gopkg.in/tylerb/graceful.v1      # graceful server shutdown
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

$GOPATH/bin/ginkgo -r --randomizeAllSpecs -p -nodes=4

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
Coveralls.io link: [https://coveralls.io/r/sogko/slumber]

To generate coverage profile

```bash
cd $GOPATH/src/github.com/sogko/slumber

# run test recursively and generate coverage data for each package
$GOPATH/bin/ginkgo -r -cover -p

# join coverage data for all packages into a single profile (for coveralls.io)
$GOPATH/bin/gover . slumber.coverprofile
```

To view coverage

```bash
go tool cover -html=$GOPATH/src/github.com/sogko/slumber/slumber.coverprofile
```

## Sessions Management
To generate key pair for signing JWT claims
```bash
$ openssl genrsa -out demo.rsa 1024 # the 1024 is the size of the key we are generating
$ openssl rsa -in demo.rsa -pubout > demo.rsa.pub 
```

## Architecture
<a href="http://i.imgur.com/HwIhPz7.png"><img src="http://i.imgur.com/HwIhPz7.png"/ height="750"/></a>


## TODO
* [x] Add open-source license (MIT) 
* [x] API versioning using Accept header, for e.g: `Accept=application/json;version=1.0,*/*`
* [x] User and roles management
* [x] Session management using JWT tokens
* [x] Activity-based access control (ABAC)
* [x] Refactor using Dependency Inversion
* [ ] Task scheduler
* [ ] Load test using vegeta
* [ ] i18n (internationalisation)
* [ ] Implement another router library for kicks
* [ ] Consolidate util libraries and publish as separate package
* [x] Abstract away negroni middlewares to a generic http.HandlerFunc
* [x] Create a REST API server project using this package as a boilerplate without changing this package
* [ ] User documentation (create resources routes, ACL, tests, middlewares etc)
