# cleanup dependencies and download missing ones
.PHONY: deps
deps:
	go mod tidy
	go mod download

# run dependency cleanup, followed by updating the patch version
.PHONY: deps-update
deps-update: deps
	go get -u=patch
	
# run tests
.PHONY: tests
tests:
	go test -race -cover -count 1 ./...

# run standard go tooling for better rcode hygiene
.PHONY: tidy
tidy: imports fmt
	go vet ./...
	golint ./...

# automatically add missing imports
.PHONY: imports
imports:
	find . -type f -name '*.go' -exec goimports -w {} \;

# format code and simplify if possible
.PHONY: fmt
fmt:
	find . -type f -name '*.go' -exec gofmt -s -w {} \;

verifiers: staticcheck

staticcheck:
	@echo "Running $@ check"
	@GO111MODULE=on ${GOPATH}/bin/staticcheck ./...


gen:
	protoc \
        -I=pb \
        -I=${GOPATH}/src \
        -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
        -I=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		pb/email.proto \
        --gogofaster_out=plugins=grpc:pb
	sed -i '/_ \"google\/protobuf\"/d' pb/email.pb.go