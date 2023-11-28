up:
	docker compose up -d

stop:
	docker compose stop

down:
	docker compose down --remove-orphans --volumes

psql:
	docker exec -it pg1 psql -U user -d test

redis:
	docker exec -it redis1 redis-cli

redis-sh:
	docker exec -it redis1 sh