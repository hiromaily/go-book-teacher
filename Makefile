# Note: tabs by space can't not used for Makefile!

###############################################################################
# Golang formatter and detection
###############################################################################
fmt:
	go fmt `go list ./... | grep -v '/vendor/'`

vet:
	go vet `go list ./... | grep -v '/vendor/'`

fix:
	go fix `go list ./... | grep -v '/vendor/'`

lint:
	golint ./... | grep -v '^vendor\/' || true
	misspell `find . -name "*.go" | grep -v '/vendor/'`
	ineffassign .

chk:
	go fmt `go list ./... | grep -v '/vendor/'`
	go vet `go list ./... | grep -v '/vendor/'`
	go fix `go list ./... | grep -v '/vendor/'`
	golint ./... | grep -v '^vendor\/' || true
	misspell `find . -name "*.go" | grep -v '/vendor/'`
	ineffassign .


###############################################################################
# Settings
###############################################################################
init:
	go get -u github.com/tools/godep


###############################################################################
# Build Local
###############################################################################
bld:
	rm -rf Godeps
	rm -rf ./vendor
	go build -i -v -o ${GOPATH}/bin/book ./cmd/book/

run:
	go run ./cmd/book/main.go
	#go run ./cmd/book/main.go -i 90
	#go run ./cmd/book/main.go -t settings.toml -i 90

godep:
	#Save
	rm -rf Godeps
	rm -rf vendor
	godep save ./...


###############################################################################
# Execution Local
###############################################################################
