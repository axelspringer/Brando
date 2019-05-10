build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/brando brando/*

test:
	go test brando/*

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean test build
	sls deploy --verbose

.PHONY: deploy-prod
deploy-prod: clean test build
	sls deploy --verbose --stage prod
