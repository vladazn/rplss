SERVICE_TAG = vladazn/rplss/service
SWAGGER_TAG = vladazn/rplss/swagger
API_TAG = vladazn/rplss/api
NGINX_TAG = vladazn/rplss/nginx
VERSION = test


all: generate swagger

generate:
	rm -rf proto/gen
	cd proto && buf generate

.PHONY: swagger
swagger:
	rm swagger/ui/swagger.json
	cp proto/gen/openapiv2/proto/game/game.swagger.json swagger/ui/swagger.json

docker_service:
	docker build -t $(SERVICE_TAG):$(VERSION) -f service/docker/Dockerfile .

docker_api:
	docker build -t $(API_TAG):$(VERSION) -f api/docker/Dockerfile .

docker_swagger:
	docker build -t $(SWAGGER_TAG):$(VERSION) -f swagger/docker/Dockerfile .

docker: docker_service docker_api docker_swagger
