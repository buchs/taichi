## Instructions

This repo serves a goal to create a demo RESTful application using Golang. It provides an employee tracking system. A member has a name and a type the late one can be an employee or a contractor - if it's a contractor, the duration of the contract needs to be saved, and if it's an employee we need to store their role, for instance: Software Engineer, Project Manager and so on. A member can be tagged, for instance: C#, Angular, General Frontend, Seasoned Leader and so on. Tags will likely be used as filters later. We need to offer a RESTful CRUD for all the information above.

## Tutorial on Use

If you are running on Linux, you should be able to build and test everything from make: `make test` will build to go binary and test it. `make docker` will take that a step further and build the docker image. The binary will be in the file **taichi**. That binary runs the server. go test in the same source directory will run the testing. Finally `make docker_test` will do the above and also test the docker image.

Tests are functional-level. 

However, if you want to do things without make, you can. 
1. Install the Go libraries listed below in section **Go Libraries Used...** into your GOPATH.
2. Build binary: `go build taichi.go tai_routes.go`
3. Indicate you want to use a test database: `export TAI_ENVIRONMENT=test`
4. Run binary server in background: `./taichi > taichi.log &`
5. Run tests: `go test`
6. Stop background server: `pkill taichi`
7. Build docker image: `docker build -t taichi:latest .`
8. Run docker image: `docker run --rm -p 3000:3000 taichi &`
9. Test running image: `go test`
10. Kill running image: `docker kill $(docker ps | grep taichi | cut -d ' ' -f 1)`

When executed in production, don't set the TAI_ENVIRONMENT variable and a production database will be used (and not erased between sessions)

The docker image is saved to the local repository with the tag: taichi:latest. 

## Design of API

The method, URL, and POST body are provided after each descriptive phrase. The following subsections enumerate the API operations that need to be handled.

### Support Create

1. create member:  **POST /create-member**  body: **name, thetype, data**
2. add a tag to existing member:  **POST /create-tag**  body: **name, tag**

### Support Read

1. get member:  **GET /read-member ?name=""**
2. get member tags:  **GET /read-tags  ?member_id=num**
3. get all members:  **GET /read-all-members**
4. find members based on tags:  **POST /find-member-tags**  body: **tags** 

### Support Update

1. update member name:  **POST /update-name** body: **member_id,new name**
2. update member type & data:  **POST /update-type**  body: **member_id, type, data**

### Support Delete

1. delete member  **DELETE /delete-member ?member_id=num**
2. delete tag from member  **DELETE /delete-member-tag ?member_id=num;tag=""**

## Data Plan:

Utilizing a sqlite3 database for this exercise. Standard relational schema is utilized.

1. Table *members* uses rowid for an integer member id, and string columns: name, thetype, and data.
2. Table *tags* uses rowid for an integer tag id, and columns: integer member id and string tag.
3. *data* will be either a name or an integer - both will be stored as a string.
4. *name* will be a string of the full name with embedded spaces, making no assumptions about number of names per person or sorting. 
5. The name may not be unique, so the member id (integer) is used for unique identification.

## Testing Forecast
These tests are written to be executed in sequence and assuming an empty database at the start.

#### Test 1 - create new member
* request: POST /create-member, body (json): {"name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
* response: 201  body (json) {"member_id": "1"}

#### Test 2 - read created member back
* request: GET /read-member?name="Kevin Buchs"
* response: 200, json = {"member_id": "1", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}

#### Test 3 - add tag to member
* request: POST /create-tag, body (json): {"member_id": "1", "tag": "C#"} 
* response: 201

#### Test 4 - add second member
* request: POST /create-member, body (json): {"name": "Adam Buchs", "thetype": "contractor", "data": "6"}
* response: 201  body (json) {"member_id": "2"}

#### Test 5 - add tag to second member
* request: POST /create-tag, body (json): {"member_id": "2", "tag": "Angular"}
* response: 201
  
#### Test 6 - get all members
* request: GET /read-all-members
* response: 200, json = [
    {"member_id": "1", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}, {"member_id": "2", "name": "Adam Buchs", "thetype": "contractor", "data": "6"}
    ]
  
#### Test 7 - find members with tag "C#"
* request: POST /find-member-tags  body (json): {"tags": ["C#"]}
* response: 200, json = [
    {"member_id": "1", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
    ]
    
#### Test 8 - find members with tag "Angular"
* request: POST /find-member-tags  body (json): {"tags": ["Angular"]}
* response: 200, json = [
     {"member_id": "2", "name": "Adam Buchs", "thetype": "contractor", "data": "6"}
    ]

#### Test 9 - find member on lable "Jugular"
* request: POST /find-member-tags  body (json): {"tags": ["Jugular"]}
* response: 200, json = {[]}

#### Test 10 - update member type adam to Employee and data to Software Engineer
* request: POST /update-type  body (json): {"member_id": "2", "thetype": "employee", "Data": "Software Engineer"}
* response: 201
  
#### Test 11 - update member name Adam to Kevin - creating name duplicates
* request: POST /update-name  body: {"member_id": "2", "name": "Kevin Buchs"}
* response: 201
  
#### Test 12 - get member "Kevin Buchs"
* request: GET /read-member?name="Kevin Buchs"
* response: 200, json = [
    {"member_id": "1", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}, {"member_id": "2", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
     ]

#### Test 13 - get all members
* request: GET /read-all-members
* response: 200, json = [
    {"member_id": "1", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"},
    {"member_id": "2", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
    ]

#### Test 14 - delete tag from member Kevin2 "Angular"
* request: DELETE /delete-member-tag?member_id=2;tag="Angular"
* response: 200

#### Test 15 - find member on tags "Angular"
* request: POST /find-member-tags  body (json): {"tags": ["Angular"]}
* response: 200, json = {[]}

#### Test 16 - delete member Kevin 1
* request: DELETE /delete-member ?member_id=1
* response: 200

#### Test 17 - find member on tags "C#"
* request: POST /find-member-tags  body (json): {"tags": ["C#"]}
* response: 200, json = {[]}

#### Test 18 - get all members
* request: GET /read-all-members
* response: 200, json = [
    {"member_id": "2", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
    ]

## Go Libraries Used - Installed By Makefile

* I make use of the redirector for Go packages at http://labix.org/gopkg.in . It allows one to fix the version of the library used. 
* What follows is the correspondence between the actual packages and the gopkg.in paths
  
    | Package Path | Version | gopkg.in Path |
    | github.com/antonholmquist/jason | v1 | gopkg.in/antonholmquist/jason.v1 |
    | github.com/go-chi/chi | v4 | gopkg.in/go-chi/chi.v4 |
    | github.com/mattn/go-sqlite3 | v1 | gopkg.in/mattn/go-sqlite3.v1 |

* Make will install these into a local Go subdirectory, which is used for the GOPATH. 

## Other Dependencies

* **sqlite3** is utilized for the database. It is built-in with the docker image. For local testing, you will need to have it installed.
* built and tested with **Go 1.12.4** and **1.12.5**.

## Concerns

1. Contract duration may not be useful without a start or end date. I suggest, instead of duration or in addition to it, we have an end date

## Potential Future Enhancements

1. Allow searching for multiple tags with and/or operators
2. Using RowID in the members table is fragile - should instead have separate member id. 
3. Of course, there is NO security here - should at least have some authentication
