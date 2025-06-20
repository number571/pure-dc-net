M=world
.PHONY: default run req
default: run 
run:
	docker-compose build
	docker-compose up
req:
	curl -X POST 'http://localhost:8081/dc' --data '$(M)'
