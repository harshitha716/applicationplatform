.PHONY: localdev clouddev clouddev-backend down clean

LOCAL_COMPOSE_CMD=docker compose --env-file .env -f docker-compose.yaml
FRONTEND_SERVICES="dashboard"
CLOUD_COMPOSE_CMD=docker compose -f docker-compose.yaml -f docker-compose.cloud.yaml --env-file .env --env-file .env.cloud


localdev:
	$(LOCAL_COMPOSE_CMD) up -d --build

localdev-backend:
	@services=$$($(LOCAL_COMPOSE_CMD) config --services | grep -vE "$(FRONTEND_SERVICES)"); \
		echo "Starting backend services: $$services"; \
		$(LOCAL_COMPOSE_CMD) up -d --build $$services

clouddev:
	$(CLOUD_COMPOSE_CMD) up -d --build

clouddev-backend:
	@services=$$($(CLOUD_COMPOSE_CMD) config --services | grep -vE "$(FRONTEND_SERVICES)"); \
		echo "Starting backend services: $$services"; \
		$(CLOUD_COMPOSE_CMD) up -d --build $$services

down:
	$(CLOUD_COMPOSE_CMD) down

clean:
	$(CLOUD_COMPOSE_CMD) down -v
