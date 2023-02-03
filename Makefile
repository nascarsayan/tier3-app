all: start

start:
	docker compose up -d

stop:
	docker compose down

# rebuild all images and restart containers
refresh:
	docker compose up -d --build --force-recreate
