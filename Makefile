.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/interactionsBin interactions/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/eventHandlerBin eventhandler/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/authCallbackBin authcallback/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/selectMenuBin selectmenu/main.go	

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy-dev: clean build
	sls deploy -s dev --verbose

deploy-prod: clean build
	sls deploy -s production --verbose
