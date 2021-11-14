runlocal: export AWS_DEFAULT_REGION=eu-west-1
runlocal: export AWS_ACCESS_KEY_ID=test
runlocal: export AWS_SECRET_ACCESS_KEY=eu-west-1
runlocal: export AWS_EC2_CUSTOM_ENDPOINT=http://localhost:4566

.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf ./artifacts 

.PHONY: artifacts
artifacts: clean
	./scripts/build.sh $(VERSION)

.PHONY: build
build:
	go build -o ./bin/ec2ti main.go

.PHONY: run
run:
	go run main.go

.PHONY: localenv
localenv:
	docker-compose up -d

.PHONY: runlocal
runlocal:
	go run main.go

.PHONY: lint
lint:
	go lint main.go
	go lint internal/*
