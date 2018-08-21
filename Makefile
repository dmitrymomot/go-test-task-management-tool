.PHONY: run down restart pma

run:
	docker-compose up -d db \
	&& sleep 10 \
	&& docker-compose up -d app \
	&& docker ps -a \
	&& docker logs testapp

stop:
	docker-compose down
	docker rmi -f testapp

restart: stop run

pma:
	docker-compose up -d pma
