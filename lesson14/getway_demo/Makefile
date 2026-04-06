COMPOSE = docker compose

.PHONY: up down logs ps test

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

test:
	go test ./...
