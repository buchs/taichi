# Targets of interest
#   taichi = default - builds the executable from Go code
#   test - runs tests
#   docker - builds docker image
#   docker_test - does all of the above and runs tests on executing container

# Behind the scenes, setup_environment installs go libraries as specified
# into a local GOPATH

GOPATH = $(CURDIR)/Go
export GOPATH

taichi : taichi.go tai_routes.go setup_environment
	go build taichi.go tai_routes.go

setup_environment :
ifeq ($(wildcard Go/pkg),)   # Won't execise this block if it finds a Go subdirectory, with pkg in it
	mkdir -p $(GOPATH)
	go get -v gopkg.in/antonholmquist/jason.v1
	go get -v gopkg.in/go-chi/chi.v4
	go get -v gopkg.in/mattn/go-sqlite3.v1
endif

test : taichi 
	TAI_ENVIRONMENT=test ./taichi > taichi.log &
	go test
	@ pkill taichi

docker : test
	docker build -t taichi:latest . 

docker_test : docker
	docker run --rm -p 3000:3000 taichi &
	go test
	docker kill $$(docker ps | grep taichi | cut -d ' ' -f 1)
