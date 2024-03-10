run:
	go run ./cmd/api/

docker-build:
	docker build -t api .

docker-run:
	docker run --network=host -p 8080:8080 api
