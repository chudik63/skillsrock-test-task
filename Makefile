build:
	docker-compose up --build

up:
	docker-compose up

down:
	docker-compose down

swagger:
	swag init -g cmd/main.go -o docs --parseDependency --parseInternal