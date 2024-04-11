run:
	docker run -d -p 5433:5432 --name postgres_banners -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admax -e POSTGRES_DB=banner postgres
	timeout 5
	go run ./cmd/banners/main.go

container:
	docker compose up