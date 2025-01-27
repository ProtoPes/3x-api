.SILENT:
.PHONY:

.DEFAULT_GOAL: build

build: clean
	go build -o ./bin/ ./cmd/api-server

run: build
	go run ./cmd/api-server $(ARG)

clean:
	rm -f ./bin/api-server

up: build
	sudo docker compose up -d --remove-orphans --build
down:
	sudo docker compose down --volumes --remove-orphans
	echo y | sudo docker image prune
