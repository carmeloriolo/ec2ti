artifacts:
	./scripts/build.sh $(VERSION)

build:
	go build -o ./bin/ec2ti cmd/main.go

run:
	go run cmd/main.go

mock:
	go run cmd/mock/main.go

lint:
	go lint cmd/main.go
	go lint internal/*
