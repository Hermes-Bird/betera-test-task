build-server:
	go build -o build/server cmd/server/main.go

build-and-run-server:
	go build -o build/server cmd/server/main.go && build/server

run-docker-app:
	docker compose build && docker compose --env-file ./.env up
