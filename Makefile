.PHONY: default run request
default: run 
run:
	docker-compose build
	docker-compose up
request:
	curl -X POST http://localhost:8081/dc --data 'world'
