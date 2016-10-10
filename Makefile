all: build golint vet master

GLOCK := glock

$(GLOCK): 
	go get github.com/robfig/glock

build: $(GLOCK)
	$(GLOCK) sync github.com/zhonglin6666/unisersal/

golint:
	golint client/
	golint cmd/...
	golint config/

	@echo "golint meta"
	@find ./meta/ -name '*.go' ! -name '*.pb.go' | xargs golint

	@echo "golint server"
	@find ./server/ -name '*.go' ! -name '*.pb.go' | xargs golint

	golint storage/

	@echo "golint storage/engine"
	@find ./storage/engine/ -name '*.go' ! -name '*.pb.go' | xargs golint

	golint pd/
	golint util/...

clean:
	@rm -rf bin

fmt:
	gofmt -s -l -w .
	goimports -l -w .

vet:
	go tool vet . 2>&1
	go tool vet --shadow . 2>&1

master:
	go build -o bin/master ./cmd/master/main.go

agent:
	go build -o bin/agent ./cmd/agent/main.go



.PHONY: master agent

test:
	go test -v ./client
	go test -v ./storage/engine
	go test -v ./storage/
	go test -v ./pd
